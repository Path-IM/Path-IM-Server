package repository

import (
	"context"
	"github.com/showurl/Path-IM-Server/common/types"
)

func (r *Rep) IncrUserSeq(uid string) (uint64, error) {
	key := types.RedisKeyUserIncrSeq + uid
	count, err := r.Cache.Incr(context.Background(), key).Result()
	return uint64(count), err
}

func (r *Rep) IncrSuperGroupSeq(groupId string) (uint64, error) {
	key := types.RedisKeyGroupIncrSeq + groupId
	count, err := r.Cache.Incr(context.Background(), key).Result()
	return uint64(count), err
}
