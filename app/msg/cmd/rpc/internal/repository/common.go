package repository

import (
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"strconv"
)

const singleGocMsgNum = 5000
const superGroupGocMsgNum = 5000

func indexGen(uid string, seqSuffix uint32) string {
	return uid + ":" + strconv.FormatInt(int64(seqSuffix), 10)
}

func getSeqUid(uid string, seq uint32) string {
	seqSuffix := seq / singleGocMsgNum
	return indexGen(uid, seqSuffix)
}
func getSeqGroupId(groupId string, seq uint32) string {
	seqSuffix := seq / superGroupGocMsgNum
	return indexGen(groupId, seqSuffix)
}

func genExceptionMessageBySeqList(seqList []uint32) (exceptionMsg []*pb.MsgData) {
	for _, v := range seqList {
		msg := new(pb.MsgData)
		msg.Seq = v
		exceptionMsg = append(exceptionMsg, msg)
	}
	return exceptionMsg
}
