package types

func GetConversationIDBySessionType(sourceID string, sessionType int) string {
	switch sessionType {
	case SingleChatType:
		return "single_" + sourceID
	case GroupChatType:
		return "supergroup_" + sourceID
	case NotificationChatType:
		return "notification_" + sourceID
	}
	return ""
}
