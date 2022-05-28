package logic

import (
	"context"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
)

func (l *MsgTransferPersistentOnlineLogic) saveSuperGroupChat(ctx context.Context, key string, c *chatpb.MsgDataToMQ) error {
	// todo 保存群聊消息
	panic("implement me")
}
