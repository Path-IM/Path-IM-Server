package pb

import "github.com/golang/protobuf/proto"

func (x *MsgData) Bytes() []byte {
	buf, _ := proto.Marshal(x)
	return buf
}
