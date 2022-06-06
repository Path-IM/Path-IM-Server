package logic

import (
	"context"
	"fmt"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/persistent/internal/repository"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/persistent/internal/svc"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/Path-IM/Path-IM-Server/common/utils"
	"github.com/Path-IM/Path-IM-Server/common/utils/statistics"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"
	"github.com/golang/protobuf/proto"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	singleMsgSuccessCount uint64
	groupMsgCount         uint64
	singleMsgFailedCount  uint64
)

func init() {
	statistics.NewStatistics(&singleMsgSuccessCount, "msg-transfer-persistent", fmt.Sprintf("%d second singleMsgCount insert to mongo", 300), 300)
	statistics.NewStatistics(&groupMsgCount, "msg-transfer-persistent", fmt.Sprintf("%d second groupMsgCount insert to mongo", 300), 300)
	statistics.NewStatistics(&groupMsgCount, "msg-transfer-persistent", fmt.Sprintf("%d second groupMsgCount insert to mongo", 300), 300)
}

type MsgTransferPersistentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	rep *repository.Rep
}

func NewMsgTransferPersistentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MsgTransferPersistentLogic {
	return &MsgTransferPersistentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		rep:    repository.NewRep(svcCtx),
	}
}

func (l *MsgTransferPersistentLogic) Do(msg []byte, msgKey string) (err error) {
	msgFromMQ := chatpb.MsgDataToMQ{}
	xtrace.StartFuncSpan(l.ctx, "MsgTransferPersistentLogic.SaveMsg.UnmarshalMsg", func(ctx context.Context) {
		err = proto.Unmarshal(msg, &msgFromMQ)
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("unmarshal msg failed, err: %v", err)
		return
	}
	logx.WithContext(l.ctx).Infof("msg: %v", msgFromMQ.String())
	isPersistent := utils.GetSwitchFromOptions(msgFromMQ.MsgData.MsgOptions, types.IsPersistent)
	switch msgFromMQ.MsgData.ConversationType {
	case types.SingleChatType:
		xtrace.StartFuncSpan(l.ctx, "MsgTransferPersistentLogic.SaveMsg.SingleChat", func(ctx context.Context) {
			if isPersistent {
				err = l.saveSingleChat(ctx, msgKey, &msgFromMQ)
				if err != nil {
					singleMsgFailedCount++
					l.Logger.Error("single data insert to mongo err ", err.Error(), " ", msgFromMQ.String())
					return
				}
				singleMsgSuccessCount++
			}
		})
	case types.GroupChatType:
		xtrace.StartFuncSpan(l.ctx, "MsgTransferPersistentLogic.SaveMsg.GroupChat", func(ctx context.Context) {
			if isPersistent {
				err = l.saveGroupChat(ctx, msgFromMQ.MsgData.ReceiveID, &msgFromMQ)
				if err != nil {
					l.Logger.Error("group data insert to mongo err ", msgFromMQ.String(), " GroupID ", msgFromMQ.MsgData.ReceiveID, " ", err.Error())
					return
				}
				groupMsgCount++
			}
		})
	default:
		l.Logger.Error("SessionType error ", msgFromMQ.String())
		return
	}
	return err
}
