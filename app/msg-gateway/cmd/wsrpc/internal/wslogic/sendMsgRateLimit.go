package wslogic

import (
	"context"
	msggatewaypb "github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
)

// 发消息频率限制
func (l *MsggatewayLogic) sendMsgRateLimit(ctx context.Context, conn *UserConn, m *msggatewaypb.Req, uid string) bool {
	takeRes, err := l.rep.SendMsgPeriodLimit.Take(uid)
	if err != nil {
		logx.WithContext(ctx).Errorf("sendMsgRateLimit take error: %v", err)
		return true
	}
	switch takeRes {
	case limit.OverQuota:
		l.sendMsgRateLimitOverQuota(ctx, conn, m, uid)
		return false
	case limit.Allowed:
		return true
	case limit.HitQuota:
		l.sendMsgRateLimitHitQuota(ctx, conn, m, uid)
		return false
	default:
		return true
	}
}

func (l *MsggatewayLogic) sendMsgRateLimitOverQuota(ctx context.Context, conn *UserConn, m *msggatewaypb.Req, uid string) {
	nReply := new(chatpb.SendMsgResp)
	nReply.ErrCode = types.ErrCodeLimit
	l.sendMsgResp(ctx, conn, m, nReply)
}

func (l *MsggatewayLogic) sendMsgRateLimitHitQuota(ctx context.Context, conn *UserConn, m *msggatewaypb.Req, uid string) {
	nReply := new(chatpb.SendMsgResp)
	nReply.ErrCode = types.ErrCodeLimit
	l.sendMsgResp(ctx, conn, m, nReply)
}
