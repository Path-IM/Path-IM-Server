package wsrepository

import (
	"fmt"
	"github.com/Path-IM/Path-IM-Server/common/types"
)

func (r *Rep) AddUserConn(uid string, platform string) error {
	key := fmt.Sprintf("%s%s", types.RedisKeyUserConn, uid)
	return r.Redis.Hset(key, platform, fmt.Sprintf("%s:%d", r.podIp, r.svcCtx.Config.RpcPort))
}

func (r *Rep) DelUserConn(uid string, platform string) error {
	key := fmt.Sprintf("%s%s", types.RedisKeyUserConn, uid)
	_, err := r.Redis.Hdel(key, platform)
	return err
}
