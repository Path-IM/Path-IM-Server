package types

// msg-gateway 用到的
const (
	WSGetNewestSeq          = 1001
	WSPullMsgBySeqList      = 1002
	WSSendMsg               = 1003
	WSGetNewestGroupSeq     = 1004
	WSPullMsgByGroupSeqList = 1005

	WSPushMsg      = 2001
	WSGroupPushMsg = 2002
)

// msg 用到的
const (
	SingleChatType       = 1
	GroupChatType        = 2
	NotificationChatType = 3

	OnlineStatus  = "online"
	OfflineStatus = "offline"
)

const (
	AtAllString = "AtAllTag"
)

// options
const (
	//OptionsKey
	IsHistory                  = "history"
	IsPersistent               = "persistent"
	IsOfflinePush              = "offlinePush"
	IsUnreadCount              = "unreadCount"
	IsConversationUpdate       = "conversationUpdate"
	IsSenderSync               = "senderSync"
	IsSenderConversationUpdate = "senderConversationUpdate"
)
