package logic

import (
	"context"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
)

func (l *MsgTransferPersistentOnlineLogic) saveNotificationChat(ctx context.Context, key string, c *chatpb.MsgDataToMQ) error {
	// todo 保存通知消息
	panic("implement me")
}
