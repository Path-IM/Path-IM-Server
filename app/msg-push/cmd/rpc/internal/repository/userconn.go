package repository

import (
	"fmt"
	"github.com/Path-IM/Path-IM-Server/common/types"
)

func (r *Rep) IsUserOnline(uid string) (bool, error) {
	key := fmt.Sprintf("%s%s", types.RedisKeyUserConn, uid)
	return r.Redis.Exists(key)
}
