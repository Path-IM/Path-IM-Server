package logic

import (
	"errors"
	chatpb "github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/showurl/Zero-IM-Server/common/types"
	"github.com/showurl/Zero-IM-Server/common/xtrace"
	"sync"
)

func (l *SendMsgLogic) sendMsgToKafka(m *chatpb.MsgDataToMQ, key string, status string) error {
	m.OperationID = xtrace.TraceIdFromContext(l.ctx)
	switch status {
	case types.OnlineStatus:
		pid, offset, err := l.svcCtx.OnlineProducer.SendMessage(l.ctx, m, key)
		if err != nil {
			l.Logger.Error(m.OperationID, "kafka send failed ", "send data ", m.String(), "pid ", pid, "offset ", offset, "err ", err.Error(), "key ", key, status)
		}
		return err
	case types.OfflineStatus:
		pid, offset, err := l.svcCtx.OfflineProducer.SendMessage(l.ctx, m, key)
		if err != nil {
			l.Logger.Error(m.OperationID, "kafka send failed ", "send data ", m.String(), "pid ", pid, "offset ", offset, "err ", err.Error(), "key ", key, status)
		}
		return err
	}
	return errors.New("status error")
}

func (l *SendMsgLogic) sendMsgToGroup(list []string, pb *chatpb.SendMsgReq, status string, sendTag *bool, wg *sync.WaitGroup) {
	//	l.Logger.Debug(pb.OperationID, "split userID ", list)
	groupPB := chatpb.SendMsgReq{Token: pb.Token, OperationID: pb.OperationID, MsgData: &chatpb.MsgData{OfflinePushInfo: &chatpb.OfflinePushInfo{}}}
	*groupPB.MsgData = *pb.MsgData
	if pb.MsgData.OfflinePushInfo != nil {
		*groupPB.MsgData.OfflinePushInfo = *pb.MsgData.OfflinePushInfo
	}
	msgToMQGroup := chatpb.MsgDataToMQ{Token: groupPB.Token, OperationID: groupPB.OperationID, MsgData: groupPB.MsgData}
	for _, v := range list {
		groupPB.MsgData.RecvID = v
		isSend := l.modifyMessageByUserMessageReceiveOpt(v, groupPB.MsgData.GroupID, types.GroupChatType, &groupPB)
		if isSend {
			msgToMQGroup.MsgData = groupPB.MsgData
			//	l.Logger.Debug(groupPB.OperationID, "sendMsgToKafka, ", v, groupID, msgToMQGroup.String())
			err := l.sendMsgToKafka(&msgToMQGroup, v, status)
			if err != nil {
				l.Logger.Error(msgToMQGroup.OperationID, "kafka send msg err:UserId", v, msgToMQGroup.String())
			} else {
				*sendTag = true
			}
		}
	}
	wg.Done()
}
