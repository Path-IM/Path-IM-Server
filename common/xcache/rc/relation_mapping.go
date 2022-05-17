package rc

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/showurl/Zero-IM-Server/common/fastjson"
	"github.com/showurl/Zero-IM-Server/common/xcache/global"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	RecordNotFoundPlaceholder float64 = -404.404 // 数据库为空的占位符
)

type RelationMapping struct {
	tx *gorm.DB
	rc redis.UniversalClient
}

func NewRelationMapping(tx *gorm.DB, rc redis.UniversalClient) *RelationMapping {
	return &RelationMapping{tx: tx, rc: rc}
}

func (r *RelationMapping) Exist(
	unknownValue string,
	model interface{},
	unknownField string,
	knownMap map[string]interface{},
	options ...Option,
) (bool, error) {
	var (
		option     = r.option(options...)
		key, e     = r.Key(model, unknownField, knownMap, options...)
		ctx        = context.Background()
		cacheExist int64
	)
	if e != nil {
		return false, e
	}
	// cacheExist
	{
		cacheExist, _ = r.rc.Exists(ctx, key).Result()
	}
	// 缓存中不存在
	if cacheExist == 0 {
		err := r.mysql2cache(ctx, key, model, unknownField, knownMap, option)
		if err != nil {
			if err == global.RedisErrorNotExists {
				return false, nil
			}
			return false, err
		}
	}
	// 取缓存中查询
	{
		_, err := r.rc.ZRank(ctx, key, unknownValue).Result()
		return !errors.Is(err, redis.Nil), nil
	}
}

type whereArgs struct {
	where string
	args  []interface{}
}

func (r *RelationMapping) mysql2cache(
	ctx context.Context,
	key string,
	model interface{},
	unknownField string,
	knownMap map[string]interface{},
	option *relationOption,
) error {
	var found []struct {
		Score  Score  `gorm:"column:sortBy;type:longtext;"`
		Member Member `gorm:"column:unknownField;type:longtext;"`
	}
	var zs []*redis.Z
	var was []whereArgs
	for k, v := range knownMap {
		if strings.Contains(k, MagicKey) {
			if strings.HasSuffix(k, gtKey) {
				was = append(was, whereArgs{
					where: fmt.Sprintf("%s > ?", strings.TrimSuffix(k, gtKey)),
					args:  []interface{}{v},
				})
				continue
			} else if strings.HasSuffix(k, ltKey) {
				was = append(was, whereArgs{
					where: fmt.Sprintf("%s < ?", strings.TrimSuffix(k, ltKey)),
					args:  []interface{}{v},
				})
				continue
			}
			if strings.HasSuffix(k, gteKey) {
				was = append(was, whereArgs{
					where: fmt.Sprintf("%s >= ?", strings.TrimSuffix(k, gteKey)),
					args:  []interface{}{v},
				})
				continue
			} else if strings.HasSuffix(k, lteKey) {
				was = append(was, whereArgs{
					where: fmt.Sprintf("%s <= ?", strings.TrimSuffix(k, lteKey)),
					args:  []interface{}{v},
				})
				continue
			}
			if strings.HasSuffix(k, whereKey) {
				if args, ok := v.([]interface{}); !ok {
					return global.ErrWhereMapValueMustInterfaceSlice
				} else {
					was = append(was, whereArgs{
						where: "(" + strings.TrimSuffix(k, whereKey) + ")",
						args:  args,
					})
					continue
				}
			}
		}
		was = append(was, whereArgs{
			where: fmt.Sprintf("%s = ?", k),
			args:  []interface{}{v},
		})
	}
	name, er := r.tableName(model)
	if er != nil {
		return er
	}
	tx := r.tx.Table(name).Select(fmt.Sprintf(
		"%s as unknownField,%s as sortBy",
		unknownField, strings.Split(option.order, " ")[0],
	))
	for _, wa := range was {
		tx = tx.Where(wa.where, wa.args...)
	}
	err := tx.
		Order(option.order).
		Limit(option.size).
		Find(&found).Error
	if err != nil {
		return err
	}
	for _, f := range found {
		if f.Member.Data != nil {
			zs = append(zs, &redis.Z{
				Score:  f.Score.Float64(),
				Member: f.Member.Data,
				//Member: reflect.ValueOf(f.Member.Data).Elem().Elem().Interface(),
			})
		}
	}
	var e error
	if len(zs) == 0 {
		// 占位符
		zs = append(zs, &redis.Z{
			Score:  RecordNotFoundPlaceholder,
			Member: "nil",
		})
		e = global.RedisErrorNotExists
	}
	err1 := r.rc.ZAdd(ctx, key, zs...).Err()
	if err1 != nil {
		return err1
	}
	var expiredTime = time.Minute
	if d, ok := model.(IRelationExpired); ok {
		expiredTime = global.ExpireDuration(d.RelationExpiredSecond())
	}
	r.rc.Expire(ctx, key, expiredTime)
	return e
}

func (r *RelationMapping) List(
	results interface{},
	start int64,
	stop int64,
	model interface{},
	unknownField string,
	knownMap map[string]interface{},
	options ...Option,
) error {
	var (
		option     = r.option(options...)
		key, e     = r.Key(model, unknownField, knownMap, options...)
		ctx        = context.Background()
		cacheExist int64
		revr       bool
	)
	if e != nil {
		return e
	}
	if strings.HasSuffix(option.order, "DESC") || strings.HasSuffix(option.order, "desc") {
		revr = true
	}
	// cacheExist
	{
		cacheExist, _ = r.rc.Exists(ctx, key).Result()
	}
	// 取缓存中的数据
	getByCache := func(results interface{}) error {
		var stopR int64
		if stop <= 0 {
			stopR = -1
		} else {
			stopR = stop
		}
		var (
			result []redis.Z
			err    error
		)
		if !revr {
			result, err = r.rc.ZRangeWithScores(ctx, key, start, stopR).Result()
		} else {
			result, err = r.rc.ZRevRangeWithScores(ctx, key, start, stopR).Result()
		}
		if err != nil {
			return err
		}
		// -404 是空占位符
		if len(result) == 1 && result[0].Score == RecordNotFoundPlaceholder {
			return nil
		}
		var rs []interface{}
		for _, res := range result {
			rs = append(rs, res.Member)
		}
		buf, _ := fastjson.Marshal(rs)
		err = fastjson.Unmarshal(buf, results)
		if err != nil {
			return err
		}
		return nil
	}
	if cacheExist != 0 {
		return getByCache(results)
	}
	// 取数据库的数据
	err := r.mysql2cache(ctx, key, model, unknownField, knownMap, option)
	if err != nil {
		return err
	}
	return getByCache(results)
}

func (r *RelationMapping) ListByScore(
	results interface{},
	by *redis.ZRangeBy,
	model interface{},
	unknownField string,
	knownMap map[string]interface{},
	options ...Option,
) error {
	var (
		option     = r.option(options...)
		key, e     = r.Key(model, unknownField, knownMap, options...)
		ctx        = context.Background()
		cacheExist int64
		revr       bool
	)
	if e != nil {
		return e
	}
	if strings.HasSuffix(option.order, "DESC") || strings.HasSuffix(option.order, "desc") {
		revr = true
	}
	// cacheExist
	{
		cacheExist, _ = r.rc.Exists(ctx, key).Result()
	}
	// 取缓存中的数据
	getByCache := func(results interface{}) error {
		var (
			result []redis.Z
			err    error
		)
		if !revr {
			// 正序
			result, err = r.rc.ZRangeByScoreWithScores(ctx, key, by).Result()
		} else {
			result, err = r.rc.ZRevRangeByScoreWithScores(ctx, key, by).Result()
		}
		if err != nil {
			return err
		}
		// -404 是空占位符
		if len(result) == 1 && result[0].Score == RecordNotFoundPlaceholder {
			return nil
		}
		var rs []interface{}
		for _, res := range result {
			rs = append(rs, res.Member)
		}
		buf, _ := fastjson.Marshal(rs)
		err = fastjson.Unmarshal(buf, results)
		if err != nil {
			return err
		}
		return nil
	}
	if cacheExist != 0 {
		return getByCache(results)
	}
	// 取数据库的数据
	err := r.mysql2cache(ctx, key, model, unknownField, knownMap, option)
	if err != nil {
		return err
	}
	return getByCache(results)
}

// 返回 []redis.Z
func (r *RelationMapping) ListByScoreWithZ(
	results interface{},
	by *redis.ZRangeBy,
	model interface{},
	unknownField string,
	knownMap map[string]interface{},
	options ...Option,
) ([]redis.Z, error) {
	var (
		option     = r.option(options...)
		key, e     = r.Key(model, unknownField, knownMap, options...)
		ctx        = context.Background()
		cacheExist int64
		revr       bool
	)
	if e != nil {
		return nil, e
	}
	if strings.HasSuffix(option.order, "DESC") || strings.HasSuffix(option.order, "desc") {
		revr = true
	}
	// cacheExist
	{
		cacheExist, _ = r.rc.Exists(ctx, key).Result()
	}
	// 取缓存中的数据
	getByCache := func(results interface{}) ([]redis.Z, error) {
		var (
			result []redis.Z
			err    error
		)
		if !revr {
			// 正序
			result, err = r.rc.ZRangeByScoreWithScores(ctx, key, by).Result()
		} else {
			result, err = r.rc.ZRevRangeByScoreWithScores(ctx, key, by).Result()
		}
		if err != nil {
			return result, err
		}
		// -404 是空占位符
		if len(result) == 1 && result[0].Score == RecordNotFoundPlaceholder {
			return nil, nil
		}
		var rs []interface{}
		for _, res := range result {
			rs = append(rs, res.Member)
		}
		buf, _ := fastjson.Marshal(rs)
		err = fastjson.Unmarshal(buf, results)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	if cacheExist != 0 {
		return getByCache(results)
	}
	// 取数据库的数据
	err := r.mysql2cache(ctx, key, model, unknownField, knownMap, option)
	if err != nil {
		return nil, err
	}
	return getByCache(results)
}

func (r *RelationMapping) tableName(model interface{}) (string, error) {
	var (
		tableName = ""
	)
	// 获取表名
	{
		if t, ok := model.(schema.Tabler); !ok {
			return "", global.ErrTablerNotImplement
		} else {
			tableName = t.TableName()
		}
	}
	return tableName, nil
}

// Key
// k = $(zset:relation:$tableName:zsize:$size:sortby:$order:$key1:$value1:$key2:$value2:unknownField:$unknownField).trim()
// k = $(z:r:$tableName:zs:$size:sb:$order:$key1:$value1:$key2:$value2:uf:$unknownField).trim()
// v =
func (r *RelationMapping) Key(
	model interface{},
	unknownField string,
	knownMap map[string]interface{},
	options ...Option,
) (string, error) {
	option := r.option(options...)
	var (
		tableName = ""
		where     = "" // key1:value1:key2:value2
	)
	// 获取表名
	{
		var err error
		tableName, err = r.tableName(model)
		if err != nil {
			return tableName, err
		}
	}
	// where
	{
		var keys []string
		for k := range knownMap {
			keys = append(keys, k)
		}
		sort.StringsAreSorted(keys)
		for _, k := range keys {
			v := knownMap[k]
			where += fmt.Sprintf("%s:%v", k, v)
		}
	}
	return strings.ReplaceAll(url.QueryEscape(fmt.Sprintf(
		"z:r:%s:zs:%d:sb:%s:%s:uf:%s",
		tableName,
		option.size,
		option.order,
		where,
		unknownField,
	)), "%3A", ":"), nil
}

func (r *RelationMapping) option(options ...Option) *relationOption {
	option := defaultRelationOption()
	for _, o := range options {
		o(option)
	}
	return option
}
