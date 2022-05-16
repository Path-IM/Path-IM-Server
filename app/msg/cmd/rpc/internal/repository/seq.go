package repository

import (
	"context"
	"github.com/showurl/Zero-IM-Server/common/types"
)

func (r *Rep) GetUserSeq(uid string) (uint64, error) {
	key := types.RedisKeyUserIncrSeq + uid
	count, err := r.Cache.Get(context.Background(), key).Int64()
	return uint64(count), err
}

func (r *Rep) GetSuperGroupSeq(groupId string) (uint64, error) {
	key := types.RedisKeySuperGroupIncrSeq + groupId
	count, err := r.Cache.Get(context.Background(), key).Int64()
	return uint64(count), err
}
