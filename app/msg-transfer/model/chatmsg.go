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

type CassUserChat struct {
	Uid  string
	Msgs []map[int64][]byte // map[sendtime]pb.MsgData
}
type CassGroupChat struct {
	Groupid string
	Msgs    []map[int64][]byte // map[sendtime]pb.MsgData
}
