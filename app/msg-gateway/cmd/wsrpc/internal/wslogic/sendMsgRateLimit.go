package wslogic

import (
	"context"
	imuserpb "github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"
	msggatewaypb "github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// 发消息频率限制
func (l *MsggatewayLogic) sendMsgRateLimit(ctx context.Context, conn *UserConn, m *msggatewaypb.BodyReq, uid string, platformID string, req *chatpb.SendMsgReq) bool {
	takeRes, err := l.rep.SendMsgPeriodLimit.TakeWithContext(ctx, uid)
	if err != nil {
		logx.WithContext(ctx).Errorf("sendMsgRateLimit take error: %v", err)
		return true
	}
	switch takeRes {
	case limit.OverQuota, limit.HitQuota:
		l.sendMsgRateLimitErr(ctx, conn, m, uid, platformID, req)
		return false
	default:
		return true
	}
}

func (l *MsggatewayLogic) sendMsgRateLimitErr(ctx context.Context, conn *UserConn, m *msggatewaypb.BodyReq, uid string, platformID string, req *chatpb.SendMsgReq) {
	if l.svcCtx.Config.EnableUserCallback {
		_, err := l.svcCtx.ImUserService().UserCallback(l.ctx, &imuserpb.UserCallbackReq{
			Event:      imuserpb.UserCallbackReq_MsgTooOften,
			Timestamp:  time.Now().UnixMilli(),
			UserID:     uid,
			Platform:   platformID,
			RemoteAddr: conn.RemoteAddr().String(),
		})
		if err != nil {
			logx.Errorf("user callback err %s", err.Error())
			err = nil
		}
	}
	nReply := new(chatpb.SendMsgResp)
	nReply.ErrCode = types.ErrCodeLimit
	nReply.ErrMsg = "请求太频繁"
	nReply.ClientMsgID = req.MsgData.ClientMsgID
	nReply.ReceiveID = req.MsgData.ReceiveID
	nReply.ConversationType = req.MsgData.ConversationType
	nReply.ContentType = req.MsgData.ContentType
	l.sendMsgResp(ctx, conn, m, nReply, uid, platformID)
}
