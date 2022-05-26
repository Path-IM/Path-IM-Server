package rpclogic

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/showurl/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcsvc"
	"github.com/showurl/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wslogic"
	"github.com/showurl/Path-IM-Server/common/types"
	"github.com/showurl/Path-IM-Server/common/xtrace"

	"github.com/showurl/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnlinePushMsgLogic struct {
	ctx    context.Context
	svcCtx *rpcsvc.ServiceContext
	logx.Logger
}

func NewOnlinePushMsgLogic(ctx context.Context, svcCtx *rpcsvc.ServiceContext) *OnlinePushMsgLogic {
	return &OnlinePushMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OnlinePushMsgLogic) OnlinePushMsg(in *pb.OnlinePushMsgReq) (*pb.OnlinePushMsgResp, error) {
	logic := wslogic.NewMsggatewayLogic(nil, nil)
	var resp []*pb.SingleMsgToUser
	msgBytes, _ := proto.Marshal(in.MsgData)
	reqIdentifier := types.WSPushMsg
	if in.MsgData.SessionType == types.GroupChatType {
		reqIdentifier = types.WSGroupPushMsg
	}
	mReply := &pb.Resp{
		ReqIdentifier: uint32(reqIdentifier),
		Data:          msgBytes,
	}
	replyBytes, err := proto.Marshal(mReply)
	if err != nil {
		l.Error("data encode err ", err.Error())
	}
	var tag bool
	recvID := in.PushToUserID
	platformList := []string{
		types.IOSPlatformStr,
		types.AndroidPlatformStr,
		types.WindowsPlatformStr,
		types.OSXPlatformStr,
		types.WebPlatformStr,
		types.MiniWebPlatformStr,
		types.LinuxPlatformStr,
	}
	for _, v := range platformList {
		if conn := logic.GetUserConn(recvID, v); conn != nil {
			tag = true
			var resultCode int64
			xtrace.StartFuncSpan(l.ctx, "OnlinePushMsgLogic.OnlinePushMsg", func(ctx context.Context) {
				resultCode = logic.SendMsgToUser(ctx, conn, replyBytes, in, v, recvID)
			})
			temp := &pb.SingleMsgToUser{
				ResultCode:     resultCode,
				RecvID:         recvID,
				RecvPlatFormID: types.PlatformNameToID(v),
			}
			resp = append(resp, temp)
		} else {
			temp := &pb.SingleMsgToUser{
				ResultCode:     -1,
				RecvID:         recvID,
				RecvPlatFormID: types.PlatformNameToID(v),
			}
			resp = append(resp, temp)
		}
	}
	if !tag {
	}
	return &pb.OnlinePushMsgResp{
		Resp: resp,
	}, nil
}
