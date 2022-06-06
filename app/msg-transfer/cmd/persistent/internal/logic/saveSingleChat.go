package logic

import (
	"context"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
)

func (l *MsgTransferPersistentLogic) saveSingleChat(ctx context.Context, key string, c *chatpb.MsgDataToMQ) error {
	// todo 保存单聊消息
	panic("implement me")
}
