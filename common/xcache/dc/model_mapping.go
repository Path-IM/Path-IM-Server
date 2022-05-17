package dc

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/showurl/Zero-IM-Server/common/fastjson"
	"github.com/showurl/Zero-IM-Server/common/utils/deepcopy"
	"github.com/showurl/Zero-IM-Server/common/xcache/global"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"time"
)

const (
	StringDetailPrefix        = "string:detail"
	RecordNotFoundPlaceholder = "@-1" // 数据库为空的占位符
)

// DbMapping 数据库模型的映射
type DbMapping struct {
	rc redis.UniversalClient
	tx *gorm.DB
}

func GetDbMapping(rc redis.UniversalClient, tx *gorm.DB) *DbMapping {
	return &DbMapping{rc: rc, tx: tx}
}

func (v *DbMapping) FirstByID(model interface{}, options ...FuncFirstOption) error {
	var (
		tablerName string // 表名
		ctx        = context.Background()
		id         string
	)
	// tableName
	{
		if tabler, ok := model.(schema.Tabler); !ok {
			// 未实现 Tabler 接口
			return global.ErrTablerNotImplement
		} else {
			tablerName = tabler.TableName()
		}
	}
	// id
	{
		if iGetID, ok := model.(global.IGetId); !ok {
			return global.ErrIGetIDNotImplement
		} else {
			id = iGetID.GetIdString()
		}
	}
	var (
		option = defaultOption(id)
	)
	{
		for _, o := range options {
			o(option)
		}
	}
	redisKey := v.GetKey(tablerName, id, options...)
	result, err := v.rc.Get(ctx, redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			var expiredTime = time.Minute
			if d, ok := model.(IDetailExpired); ok {
				expiredTime = global.ExpireDuration(d.DetailExpiredSecond())
			}
			// 去数据库里查询
			err = v.tx.Model(model).Where(option.where, option.args...).First(model).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					v.rc.Set(ctx, redisKey, RecordNotFoundPlaceholder, expiredTime)
				}
				return err
			} else {
				// 查询成功
				buf, err2 := fastjson.Marshal(model)
				if err2 != nil {
					return err2
				}
				v.rc.Set(ctx, redisKey, string(buf), expiredTime)
				return nil
			}
		}
		return err
	}
	if result == RecordNotFoundPlaceholder {
		return redis.Nil
	}
	err = fastjson.Unmarshal([]byte(result), model)
	if err != nil {
		return err
	}
	return nil
}

func (v *DbMapping) ListByID(
	model,
	list interface{},
	ids []string,
	options ...FuncFirstOption,
) error {
	if len(ids) == 0 {
		return nil
	}
	var (
		ctx        = context.Background()
		keys       []string
		tablerName string
		cacheIdMap = make(map[string]bool)
		needDbIds  []string
		option     = defaultListOption()
	)
	// tablerName
	{
		if tabler, ok := model.(schema.Tabler); !ok {
			return global.ErrTablerNotImplement
		} else {
			tablerName = tabler.TableName()
		}
	}
	// id
	{
		if _, ok := model.(global.IGetId); !ok {
			return global.ErrIGetIDNotImplement
		}
	}
	{
		for _, listOption := range options {
			listOption(option)
		}
	}
	// keys
	{
		for _, id := range ids {
			keys = append(keys, v.GetKey(tablerName, id, options...))
		}
	}
	result, err := v.rc.MGet(ctx, keys...).Result()
	if err != nil {
		return err
	}
	objT := reflect.TypeOf(list)
	objV := reflect.ValueOf(list)
	if objT.Kind() == reflect.Ptr {
		objT = objT.Elem()
		objV = objV.Elem()
	} else {
		return global.ErrInputListNotPtr
	}
	if objT.Kind() == reflect.Slice {
		sliceElem := objT.Elem()
		if sliceElem.Kind() == reflect.Ptr {
			sliceElem = sliceElem.Elem()
		} else {
			return global.ErrInputModelNotPtr
		}
		if sliceElem.Kind() == reflect.Struct {
			for _, i := range result {
				copyModel := deepcopy.Copy(model)
				if s, ok := i.(string); ok {
					if s == RecordNotFoundPlaceholder {
						continue
					}
					buf := []byte(s)
					err = fastjson.Unmarshal(buf, copyModel)
					if err != nil {
						continue
					}
					id := copyModel.(global.IGetId).GetIdString()
					cacheIdMap[id] = true
					objV.Set(reflect.Append(objV, reflect.ValueOf(copyModel)))
				}
			}
		} else {
			return global.ErrInputModelNotStruct
		}
	} else {
		return global.ErrInputListNotSlice
	}
	// 对比缓存中的和需要查询的
	idsFmt := ","
	{
		for _, id := range ids {
			idsFmt += id + ","
			if _, exist := cacheIdMap[id]; !exist {
				needDbIds = append(needDbIds, id)
			}
		}
	}
	var expiredTime = time.Minute
	if d, ok := model.(IDetailExpired); ok {
		expiredTime = global.ExpireDuration(d.DetailExpiredSecond())
	}
	if len(needDbIds) > 0 {
		dbList := deepcopy.Copy(list)
		//dbList := reflect.MakeSlice(objT, 0, len(needDbIds))
		//reflect.utils.Copy(dbList, objV)
		orderExpr := fmt.Sprintf("INSTR('%s',CONCAT(',',%s,','))", idsFmt, option.fieldId)
		tx := v.tx.Table(tablerName).
			Where(option.fieldId+" in (?)", needDbIds)
		if option.where != "" {
			tx = tx.Where(option.where, option.args...)
		}
		err = tx.
			Order(orderExpr).Find(dbList).Error
		if err != nil {
			return err
		}
		objV1 := reflect.ValueOf(dbList).Elem()
		for i := 0; i < objV1.Len(); i++ {
			m := objV1.Index(i).Interface()
			buf, err := fastjson.Marshal(m)
			if err != nil {
				continue
			}
			id := m.(global.IGetId).GetIdString()
			v.rc.Set(ctx, v.GetKey(tablerName, id, options...), string(buf), expiredTime)
			cacheIdMap[id] = true
			objV.Set(reflect.Append(objV, reflect.ValueOf(m)))
		}
	}
	// 对比缓存中的和需要查询的
	{
		needDbIds = []string{}
		for _, id := range ids {
			if _, exist := cacheIdMap[id]; !exist {
				needDbIds = append(needDbIds, id)
			}
		}
	}
	// 缓存制空
	{
		for _, id := range needDbIds {
			v.rc.Set(ctx, v.GetKey(tablerName, id, options...), RecordNotFoundPlaceholder, expiredTime)
		}
	}
	return err
}

func (v *DbMapping) FlushByID(model interface{}, ids []string, options ...FuncFirstOption) error {
	var (
		ctx = context.Background()
	)
	keys, err := v.GetKeys(model, ids, options...)
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return v.rc.Del(ctx, keys...).Err()
	}
	return nil
}

func (v *DbMapping) MapByID(
	model,
	list interface{},
	mp interface{},
	ids []string,
	options ...FuncFirstOption,
) error {
	if len(ids) == 0 {
		return nil
	}
	// id
	{
		if _, ok := model.(global.IGetId); !ok {
			return global.ErrIGetIDNotImplement
		}
	}
	objT := reflect.TypeOf(mp)
	objV := reflect.ValueOf(mp)
	kind := objT.Kind()
	if kind != reflect.Map {
		return global.ErrInputMpNotMap
	}
	err := v.ListByID(model, list, ids, options...)
	if err != nil {
		return err
	}
	listElem := reflect.ValueOf(list).Elem()
	for i := 0; i < listElem.Len(); i++ {
		item := listElem.Index(i).Interface()
		key := item.(global.IGetId).GetIdString()
		keyValue := reflect.ValueOf(key)
		value := reflect.ValueOf(item)
		objV.SetMapIndex(keyValue, value)
	}
	return nil
}

func (v *DbMapping) GetKeys(model interface{}, ids []string, options ...FuncFirstOption) (keys []string, err error) {
	var (
		ctx        = context.Background()
		tablerName string
		option     = defaultListOption()
	)
	// tablerName
	{
		if tabler, ok := model.(schema.Tabler); !ok {
			err = global.ErrTablerNotImplement
			return
		} else {
			tablerName = tabler.TableName()
		}
	}
	{
		for _, firstOption := range options {
			firstOption(option)
		}
	}
	if len(ids) == 0 {
		keys, err = v.rc.Keys(ctx, v.GetKey(tablerName, "*", options...)).Result()
		if err != nil {
			return
		}
	} else {
		for _, id := range ids {
			keys = append(keys, v.GetKey(tablerName, id, options...))
		}
	}
	return
}

func (v *DbMapping) GetKeyByModel(model interface{}, id string, options ...FuncFirstOption) (string, error) {
	var (
		tableName string
		option    = defaultListOption()
	)
	{
		for _, firstOption := range options {
			firstOption(option)
		}
	}
	// tablerName
	{
		if tabler, ok := model.(schema.Tabler); !ok {
			return "", global.ErrTablerNotImplement
		} else {
			tableName = tabler.TableName()
		}
	}
	return v.GetKey(tableName, id, options...), nil
}

func (v *DbMapping) GetKey(tableName string, id string, options ...FuncFirstOption) string {
	var option = defaultListOption()
	{
		for _, firstOption := range options {
			firstOption(option)
		}
	}
	return global.MergeKey(StringDetailPrefix, tableName, id) + option.keySuffix
}
