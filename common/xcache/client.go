package xcache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/showurl/Zero-IM-Server/common/xcache/global"
	zeroredis "github.com/zeromicro/go-zero/core/stores/redis"
	"log"
)

func GetClient(
	cfg zeroredis.RedisConf,
	db global.DB,
) redis.UniversalClient {
	// 打印配置
	log.Printf("redis config: %+v", cfg)
	opts := &redis.UniversalOptions{
		Addrs: []string{cfg.Host},
		DB:    db.Int(),
		//PoolSize:     15,
		//MinIdleConns: 5, // redis连接池最小空闲连接数.
		Password: cfg.Pass,
		//ReadTimeout:  5,
	}
	rc := redis.NewUniversalClient(opts)
	err := rc.Ping(context.Background()).Err()
	if err != nil {
		log.Printf("redis ping error: %+v", err)
		panic(err)
	}
	return rc
}
