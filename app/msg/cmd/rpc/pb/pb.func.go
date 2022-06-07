package pb

import (
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/golang/protobuf/proto"
)

func (x *MsgData) Bytes() []byte {
	buf, _ := proto.Marshal(x)
	return buf
}

func SetSwitchFromOptions(options *MsgOptions, key string, value bool) {
	if options == nil {
		options = &MsgOptions{}
	}
	switch key {
	case types.IsHistory:
		options.History = value
	case types.IsPersistent:
		options.Persistent = value
	case types.IsLocal:
		options.Local = value
	case types.UnreadCount:
		options.UpdateUnreadCount = value
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

func GetSwitchFromOptions(options *MsgOptions, key string, defaultValue bool) (result bool) {
	if options == nil {
		return defaultValue
	}
	switch key {
	case types.IsHistory:
		result = options.History
	case types.IsPersistent:
		result = options.Persistent
	case types.IsLocal:
		result = options.Local
	case types.UnreadCount:
		result = options.UpdateUnreadCount
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
