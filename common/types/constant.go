package types

// msg-gateway 用到的
const (
	WSGetNewestSeq          = 1001
	WSPullMsgBySeqList      = 1002
	WSGetNewestGroupSeq     = 1003
	WSPullMsgByGroupSeqList = 1004

	WSSendMsg      = 2001
	WSPushMsg      = 2002
	WSGroupPushMsg = 2003
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
