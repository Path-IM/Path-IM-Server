package utils

import (
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
)

func SetSwitchFromOptions(options *pb.MsgOptions, key string, value bool) {
	if options == nil {
		options = &pb.MsgOptions{}
	}
	switch key {
	case types.IsHistory:
		options.History = value
	case types.IsPersistent:
		options.Persistent = value
	case types.UnreadCount:
		options.UnreadCount = value
	case types.UpdateConversation:
		options.UpdateConversation = value
	case types.NeedBeFriend:
		options.NeedBeFriend = value
	case types.IsOfflinePush:
		options.OfflinePush = value
	case types.IsSenderSync:
		options.SenderSync = value
	}
}

func GetSwitchFromOptions(options *pb.MsgOptions, key string) (result bool) {
	switch key {
	case types.IsHistory:
		result = options.History
	case types.IsPersistent:
		result = options.Persistent
	case types.UnreadCount:
		result = options.UnreadCount
	case types.UpdateConversation:
		result = options.UpdateConversation
	case types.NeedBeFriend:
		result = options.NeedBeFriend
	case types.IsOfflinePush:
		result = options.OfflinePush
	case types.IsSenderSync:
		result = options.SenderSync
	}
	return
}
