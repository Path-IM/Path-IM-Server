package rpclogic

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/rpcsvc"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wslogic"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	pushpb "github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/pb"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"
	"github.com/golang/protobuf/proto"
	"github.com/zeromicro/go-zero/core/logx"
)

type PushMsgLogic struct {
	ctx    context.Context
	svcCtx *rpcsvc.ServiceContext
	logx.Logger
}

func NewPushMsgLogic(ctx context.Context, svcCtx *rpcsvc.ServiceContext) *PushMsgLogic {
	return &PushMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *PushMsgLogic) PushMsg(in *chatpb.PushMsgDataToMQ) error {
	// 判断有没有这个user
	logic := wslogic.NewMsggatewayLogic(nil, nil)
	platformList := []string{
		types.IOSPlatformStr,
		types.AndroidPlatformStr,
		types.WindowsPlatformStr,
		types.OSXPlatformStr,
		types.WebPlatformStr,
		types.MiniWebPlatformStr,
		types.LinuxPlatformStr,
	}
	reqIdentifier := types.WSPushMsg
	if in.MsgData.ConversationType == types.GroupChatType {
		reqIdentifier = types.WSGroupPushMsg
	}
	msgBytes, _ := proto.Marshal(in.MsgData)
	mReply := &pb.BodyResp{
		ReqIdentifier: uint32(reqIdentifier),
		Data:          msgBytes,
	}
	replyBytes, _ := proto.Marshal(mReply)
	var allPlatformFailed = true
	onlineUserMap := logic.GetOnlineUserMap()
	for _, recvID := range in.PushToUserID {
		if _, online := onlineUserMap[recvID]; !online {
			continue
		}
		go func(recvID string) {
			for _, v := range platformList {
				if conn := logic.GetUserConn(recvID, v); conn != nil {
					var err error
					xtrace.StartFuncSpan(l.ctx, "OnlinePushMsgLogic.OnlinePushMsg", func(ctx context.Context) {
						err = logic.SendMsgToUser(ctx, conn, replyBytes, v, recvID)
					})
					if err == nil {
						allPlatformFailed = false
					}
				}
			}
		}(recvID)
		go func(recvID string) {
			if allPlatformFailed && in.MsgData.ConversationType == types.SingleChatType {
				req := pushpb.PushMsgReq{
					MsgData:      in.MsgData,
					PushToUserID: []string{recvID},
				}
				_, err := l.svcCtx.MsgPushRpc().PushMsg(l.ctx, &req)
				if err != nil {
					l.Errorf("offline PushMsg err: %v", err)
				}
			}
		}(recvID)
	}
	return nil
}
