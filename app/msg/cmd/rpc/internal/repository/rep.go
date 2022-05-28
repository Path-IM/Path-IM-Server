package repository

import (
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/common/xcache"
	"github.com/Path-IM/Path-IM-Server/common/xcache/global"
	"github.com/Path-IM/Path-IM-Server/common/xcql"
	"github.com/Path-IM/Path-IM-Server/common/xmgo"
	"github.com/go-redis/redis/v8"
)

type Rep struct {
	svcCtx *svc.ServiceContext
	Cache  redis.UniversalClient
	IPullHistoryMsg
}

var rep *Rep

func NewRep(svcCtx *svc.ServiceContext) *Rep {
	if rep != nil {
		return rep
	}
	rep = &Rep{
		svcCtx: svcCtx,
		Cache:  xcache.GetClient(svcCtx.Config.RedisConfig.Conf, global.DB(svcCtx.Config.RedisConfig.DB)),
	}
	if svcCtx.Config.HistoryDBType == "mongo" {
		rep.IPullHistoryMsg = &MongoHistory{
			svcCtx:      svcCtx,
			MongoClient: xmgo.GetClient(svcCtx.Config.Mongo.MongoConfig),
		}
	} else if svcCtx.Config.HistoryDBType == "cassandra" {
		rep.IPullHistoryMsg = &CassandraHistory{
			svcCtx:          svcCtx,
			CassandraClient: xcql.GetClient(svcCtx.Config.Cassandra.CassandraConfig),
		}
	} else {
		panic("history db type error, select mongo or cassandra")
	}
	return rep
}
