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

	WSErrorMsg = 3001
)

// msg 用到的
const (
	SingleChatType = 1
	GroupChatType  = 2
)

const (
	AtAllString = "AtAllTag"
)

// options
const (
	//OptionsKey
	IsHistory          = "IsHistory"          // 存储到 mongodb/cassandra
	IsPersistent       = "IsPersistent"       // 存储到 持久层
	UnreadCount        = "UnreadCount"        // 是否更新未读数
	UpdateConversation = "UpdateConversation" // 更新会话
	NeedBeFriend       = "NeedBeFriend"       // 是否需要成为好友
	IsOfflinePush      = "IsOfflinePush"      // 是否离线推送
	IsSenderSync       = "IsSenderSync"       // 是否需要同步给发送者
)
