package logic

import (
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
)

func (l *SendMsgLogic) encapsulateMsgData(msg *chatpb.MsgData) {
	switch msg.ContentType {
	// todo modify options by msg.ContentType
	default:
		//utils.SetSwitchFromOptions(msg.Options, types.NeedBeFriend, true)
		chatpb.SetSwitchFromOptions(msg.MsgOptions, types.IsOfflinePush, msg.OfflinePush != nil)
		chatpb.SetSwitchFromOptions(msg.MsgOptions, types.IsSenderSync, true)
	}
}
