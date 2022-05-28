package wsrepository

import (
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wssvc"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Rep struct {
	svcCtx             *wssvc.ServiceContext
	Redis              *redis.Redis
	SendMsgPeriodLimit *limit.PeriodLimit
}

var rep *Rep

func NewRep(svcCtx *wssvc.ServiceContext) *Rep {
	if rep != nil {
		return rep
	}
	rep = &Rep{
		svcCtx: svcCtx,
		Redis:  newRedis(svcCtx.Config.Redis.Host, svcCtx.Config.Redis.Pass, svcCtx.Config.Redis.Type, svcCtx.Config.Redis.Tls),
	}
	rep.SendMsgPeriodLimit = limit.NewPeriodLimit(
		svcCtx.Config.SendMsgRateLimit.Seconds, svcCtx.Config.SendMsgRateLimit.Quota,
		rep.Redis, "periodlimit:sendmsg:", limit.Align(),
	)
	return rep
}

func newRedis(addr string, password string, typ string, tls bool) *redis.Redis {
	ops := make([]redis.Option, 0)
	if password != "" {
		ops = append(ops, redis.WithPass(password))
	}
	if typ == "cluster" {
		ops = append(ops, redis.Cluster())
	}
	if tls {
		ops = append(ops, redis.WithTLS())
	}
	return redis.New(addr, ops...)
}
