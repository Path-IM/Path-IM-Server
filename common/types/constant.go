package types

// msg-gateway 用到的
const (
	WSGetNewestSeq               = 1001
	WSPullMsgBySeqList           = 1002
	WSSendMsg                    = 1003
	WSSendSignalMsg              = 1004
	WSGetNewestSuperGroupSeq     = 1005
	WSPullMsgBySuperGroupSeqList = 1006
	WSPushMsg                    = 2001
	WSKickOnlineMsg              = 2002
	WsLogoutMsg                  = 2003
	WSSuperGroupPushMsg          = 2004
)

// msg 用到的
const (
	SingleChatType       = 1
	GroupChatType        = 2
	SuperGroupChatType   = 3
	NotificationChatType = 4

	OnlineStatus  = "online"
	OfflineStatus = "offline"
	Registered    = "registered"
	UnRegistered  = "unregistered"

	// 群聊的
	GroupCreatedNotification             = 1501
	GroupInfoSetNotification             = 1502
	JoinGroupApplicationNotification     = 1503
	MemberQuitNotification               = 1504
	GroupApplicationAcceptedNotification = 1505
	GroupApplicationRejectedNotification = 1506
	GroupOwnerTransferredNotification    = 1507
	MemberKickedNotification             = 1508
	MemberInvitedNotification            = 1509
	MemberEnterNotification              = 1510
	GroupDismissedNotification           = 1511
	GroupMemberMutedNotification         = 1512
	GroupMemberCancelMutedNotification   = 1513
	GroupMutedNotification               = 1514
	GroupCancelMutedNotification         = 1515
	GroupMemberInfoSetNotification       = 1516
)

const (
	AtAllString = "AtAllTag"
	AtNormal    = 0
	AtMe        = 1
	AtAll       = 2
	AtAllAtMe   = 3
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
	IsNotPrivate               = "notPrivate"
	IsSenderConversationUpdate = "senderConversationUpdate"
)
