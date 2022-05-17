package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/internal/svc"
	"github.com/showurl/Zero-IM-Server/common/xcache"
	"github.com/showurl/Zero-IM-Server/common/xcache/global"
	"github.com/showurl/Zero-IM-Server/common/xmgo"
	"go.mongodb.org/mongo-driver/mongo"
)

type Rep struct {
	svcCtx *svc.ServiceContext
	Cache  redis.UniversalClient
	//mysql       *gorm.DB
	MongoClient *mongo.Client
}

var rep *Rep

func NewRep(svcCtx *svc.ServiceContext) *Rep {
	if rep != nil {
		return rep
	}
	rep = &Rep{
		svcCtx: svcCtx,
		Cache:  xcache.GetClient(svcCtx.Config.RedisConfig.Conf, global.DB(svcCtx.Config.RedisConfig.DB)),
		//mysql:       xorm.GetClient(svcCtx.Config.Mysql),
		MongoClient: xmgo.GetClient(svcCtx.Config.Mongo.MongoConfig),
	}
	return rep
}
