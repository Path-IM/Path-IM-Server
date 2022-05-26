package model

type MsgInfo struct {
	SendTime int64
	Msg      []byte
}

type UserChat struct {
	UID string
	Msg []MsgInfo
}

type GroupChat struct {
	GroupID string `bson:"groupid"`
	Msg     []MsgInfo
}
