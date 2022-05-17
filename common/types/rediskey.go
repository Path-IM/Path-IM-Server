package types

const (
	RedisKeyAccountTempCode               = "ACCOUNT_TEMP_CODE"
	RedisKeyResetPwdTempCode              = "RESET_PWD_TEMP_CODE"
	RedisKeyUserIncrSeq                   = "REDIS_USER_INCR_SEQ:"        // user incr seq
	RedisKeySuperGroupIncrSeq             = "REDIS_SUPER_GROUP_INCR_SEQ:" // user incr seq
	RedisKeySuperGroupMinSeq              = "REDIS_SUPER_GROUP_MIN_SEQ:"  // user incr seq
	RedisKeyAppleDeviceToken              = "DEVICE_TOKEN"
	RedisKeyUserMinSeq                    = "REDIS_USER_MIN_SEQ:"
	RedisKeyUidPidToken                   = "UID_PID_TOKEN_STATUS:"
	RedisKeyConversationReceiveMessageOpt = "CON_RECV_MSG_OPT:"
	RedisKeyGetuiToken                    = "GETUI"
	RedisKeyUserInfoCache                 = "USER_INFO_CACHE:"
	RedisKeyFriendRelationCache           = "FRIEND_RELATION_CACHE:"
	RedisKeyBlackListCache                = "BLACK_LIST_CACHE:"
	RedisKeyGroupCache                    = "GROUP_CACHE:"
)
