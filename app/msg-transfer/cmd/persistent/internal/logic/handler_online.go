package logic

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/showurl/Zero-IM-Server/app/msg-transfer/cmd/persistent/internal/repository"
	"github.com/showurl/Zero-IM-Server/app/msg-transfer/cmd/persistent/internal/svc"
	chatpb "github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/showurl/Zero-IM-Server/common/types"
	"github.com/showurl/Zero-IM-Server/common/utils"
	"github.com/showurl/Zero-IM-Server/common/utils/statistics"
	"github.com/showurl/Zero-IM-Server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	singleMsgSuccessCount uint64
	groupMsgCount         uint64
	superGroupMsgCount    uint64
	singleMsgFailedCount  uint64
)

func init() {
	statistics.NewStatistics(&singleMsgSuccessCount, "msg-transfer-persistent", fmt.Sprintf("%d second singleMsgCount insert to mongo", 300), 300)
	statistics.NewStatistics(&groupMsgCount, "msg-transfer-persistent", fmt.Sprintf("%d second groupMsgCount insert to mongo", 300), 300)
	statistics.NewStatistics(&superGroupMsgCount, "msg-transfer-persistent", fmt.Sprintf("%d second superGroupMsgCount insert to mongo", 300), 300)
}

type MsgTransferPersistentOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	rep *repository.Rep
}

func NewMsgTransferPersistentOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MsgTransferPersistentOnlineLogic {
	return &MsgTransferPersistentOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		rep:    repository.NewRep(svcCtx),
	}
}

func (l *MsgTransferPersistentOnlineLogic) Do(msg []byte, msgKey string) {
	msgFromMQ := chatpb.MsgDataToMQ{}
	var err error
	xtrace.StartFuncSpan(l.ctx, "MsgTransferPersistentOnlineLogic.ChatMs2Mongo.UnmarshalMsg", func(ctx context.Context) {
		err = proto.Unmarshal(msg, &msgFromMQ)
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("unmarshal msg failed, err: %v", err)
		return
	}
	logx.WithContext(l.ctx).Infof("msg: %v", msgFromMQ.String())
	isPersistent := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, types.IsPersistent)
	switch msgFromMQ.MsgData.SessionType {
	case types.SingleChatType:
		xtrace.StartFuncSpan(l.ctx, "MsgTransferPersistentOnlineLogic.ChatMs2Mongo.SingleChat", func(ctx context.Context) {
			if isPersistent {
				err := l.saveSingleChat(ctx, msgKey, &msgFromMQ)
				if err != nil {
					singleMsgFailedCount++
					l.Logger.Error("single data insert to mongo err ", err.Error(), " ", msgFromMQ.String())
					return
				}
				singleMsgSuccessCount++
			}
		})
	case types.GroupChatType:
		xtrace.StartFuncSpan(l.ctx, "MsgTransferPersistentOnlineLogic.ChatMs2Mongo.SuperGroupChat", func(ctx context.Context) {
			if isPersistent {
				err := l.saveSuperGroupChat(ctx, msgFromMQ.MsgData.GroupID, &msgFromMQ)
				if err != nil {
					l.Logger.Error("super group data insert to mongo err ", msgFromMQ.String(), " GroupID ", msgFromMQ.MsgData.GroupID, " ", err.Error())
					return
				}
				superGroupMsgCount++
			}
		})
	case types.NotificationChatType:
		xtrace.StartFuncSpan(l.ctx, "MsgTransferPersistentOnlineLogic.ChatMs2Mongo.NotificationChat", func(ctx context.Context) {
			if isPersistent {
				err := l.saveNotificationChat(ctx, msgKey, &msgFromMQ)
				if err != nil {
					l.Logger.Error("single data insert to mongo err ", err.Error(), " ", msgFromMQ.String())
					return
				}
			}
		})
	default:
		l.Logger.Error("SessionType error ", msgFromMQ.String())
		return
	}
}
