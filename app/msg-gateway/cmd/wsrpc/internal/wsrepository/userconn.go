package wsrepository

import (
	"context"
	"fmt"
	"github.com/Path-IM/Path-IM-Server/common/types"
)

func (r *Rep) AddUserConn(uid string, platform string) error {
	key := fmt.Sprintf("%s%s", types.RedisKeyUserConn, uid)
	return r.Redis.Hset(key, platform, fmt.Sprintf("%s:%d", r.podIp, r.svcCtx.Config.RpcPort))
}

func (r *Rep) DelUserConn(ctx context.Context, uid string, platform string) error {
	key := fmt.Sprintf("%s%s", types.RedisKeyUserConn, uid)
	_, err := r.Redis.HdelCtx(ctx, key, platform)
	return err
}

func (r *Rep) GetUserConn(ctx context.Context, uid string) (map[string]string, error) {
	key := fmt.Sprintf("%s%s", types.RedisKeyUserConn, uid)
	return r.Redis.HgetallCtx(ctx, key)
}
