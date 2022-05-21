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
	"net/url"
	"reflect"
	"sort"
	"strings"
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

type filter struct {
	Where string
	Args  []interface{}
}

func whereMap2Filters(whereMap map[string][]interface{}) (filters []filter) {
	for where, args := range whereMap {
		filters = append(filters, filter{Where: where, Args: args})
	}
	return
}
func (v *DbMapping) FirstByWhere(model interface{}, whereMap map[string][]interface{}, options ...FuncFirstOption) error {
	var (
		tablerName string // 表名
		ctx        = context.Background()
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
	var (
		option = defaultOption()
	)
	{
		for _, o := range options {
			o(option)
		}
	}
	redisKey := v.Key(model, whereMap, options...)
	result, err := v.rc.Get(ctx, redisKey).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			var expiredTime = time.Minute
			if d, ok := model.(IDetailExpired); ok {
				expiredTime = global.ExpireDuration(d.DetailExpiredSecond())
			}
			// 去数据库里查询
			filters := whereMap2Filters(whereMap)
			tx := v.tx.Model(model).Table(tablerName)
			for _, f := range filters {
				tx = tx.Where(f.Where, f.Args...)
			}
			err = tx.Order(option.order).First(model).Error
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
func (v *DbMapping) FirstById(model interface{}, id string, options ...FuncFirstOption) error {
	return v.FirstByWhere(model, map[string][]interface{}{"id = ?": {id}}, options...)
}
func (v *DbMapping) FlushByWhere(model interface{}, whereMap map[string][]interface{}, options ...FuncFirstOption) error {
	k := v.Key(model, whereMap, options...)
	return v.rc.Del(context.Background(), k).Err()
}
func (v *DbMapping) FlushByIds(model interface{}, ids []string, options ...FuncFirstOption) error {
	var keys []string
	for _, id := range ids {
		k := v.Key(model, map[string][]interface{}{"id = ?": {id}}, options...)
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return nil
	}
	return v.rc.Del(context.Background(), keys...).Err()
}

// ListByIds
func (v *DbMapping) ListByIds(model interface{}, results interface{}, ids []string, options ...FuncFirstOption) error {
	var (
		keys       []string
		cacheIdMap = make(map[string]interface{})
	)
	for _, id := range ids {
		k := v.Key(model, map[string][]interface{}{"id = ?": {id}}, options...)
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return nil
	}
	result, err := v.rc.MGet(context.Background(), keys...).Result()
	if err != nil {
		return err
	}
	resultsV := reflect.ValueOf(results)
	list := deepcopy.Copy(results)
	objT := reflect.TypeOf(list)
	objV := reflect.ValueOf(list)
	if objT.Kind() == reflect.Ptr {
		objT = objT.Elem()
		objV = objV.Elem()
	} else {
		panic(global.ErrInputListNotPtr)
	}
	if objT.Kind() != reflect.Slice {
		panic(global.ErrInputListNotSlice)
	}
	sliceElem := objT.Elem()
	if sliceElem.Kind() != reflect.Ptr {
		panic(global.ErrInputListNotPtr)
	}
	sliceElem = sliceElem.Elem()
	if sliceElem.Kind() != reflect.Struct {
		panic(global.ErrInputModelNotStruct)
	}
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
			cacheIdMap[id] = copyModel
			objV.Set(reflect.Append(objV, reflect.ValueOf(copyModel)))
		}
	}
	for _, id := range ids {
		if value, ok := cacheIdMap[id]; !ok {
			copyModel := deepcopy.Copy(model)
			err = v.FirstById(copyModel, id, options...)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, global.RedisErrorNotExists) {
					continue
				}
				return err
			}
			resultsV.Set(reflect.Append(resultsV, reflect.ValueOf(copyModel)))
		} else {
			resultsV.Set(reflect.Append(resultsV, reflect.ValueOf(value)))
		}
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
	err := v.ListByIds(model, list, ids, options...)
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

func (v *DbMapping) Key(model interface{}, whereMap map[string][]interface{}, options ...FuncFirstOption) string {
	option := v.option(options...)
	var (
		tableName = model.(schema.Tabler).TableName()
		where     = "" // key1:value1:key2:value2
	)
	// where
	{
		var keys []string
		for k := range whereMap {
			keys = append(keys, strings.TrimSpace(k))
		}
		sort.StringsAreSorted(keys)
		for _, k := range keys {
			v := whereMap[k]
			where += fmt.Sprintf("%s:%v", k, v)
		}
	}
	return strings.ReplaceAll(url.QueryEscape(fmt.Sprintf(`%s:tablename:%s:where:%s:orderby:%s`, StringDetailPrefix, tableName, where, option.order)), "%3A", ":")
}
func (v *DbMapping) KeyById(model interface{}, id string, options ...FuncFirstOption) string {
	return v.Key(model, map[string][]interface{}{
		"id = ?": {id},
	}, options...)
}

func (v *DbMapping) option(options ...FuncFirstOption) *FirstOption {
	option := defaultOption()
	for _, o := range options {
		o(option)
	}
	return option
}
