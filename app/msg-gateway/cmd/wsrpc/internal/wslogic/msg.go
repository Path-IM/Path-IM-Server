package wslogic

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/chat"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

func (l *MsggatewayLogic) msgParse(ctx context.Context, conn *UserConn, binaryMsg []byte, uid string, platformID string) {
	m := &pb.BodyReq{}
	err := proto.Unmarshal(binaryMsg, m)
	if err != nil {
		l.sendErrMsg(ctx, conn, types.ErrCodeParams, err.Error(), uid, platformID)
		err = conn.Close()
		if err != nil {
			logx.WithContext(ctx).Error("ws close err", err.Error())
		}
		return
	}
	if err := validate.Struct(m); err != nil {
		logx.WithContext(ctx).Error("ws args validate  err", err.Error())
		l.sendErrMsg(ctx, conn, types.ErrCodeParams, err.Error(), uid, platformID)
		return
	}
	switch m.ReqIdentifier {
	case types.WSGetNewestSeq:
		l.getSeqReq(ctx, conn, m, uid, platformID)
	case types.WSGetNewestGroupSeq:
		l.getGroupSeqReq(ctx, conn, m, uid, platformID)
	case types.WSSendMsg:
		l.sendMsgReq(ctx, conn, m, uid, platformID)
	case types.WSPullMsgBySeqList:
		l.pullMsgBySeqListReq(ctx, conn, m, uid, platformID)
	case types.WSPullMsgByGroupSeqList:
		l.pullMsgByGroupSeqListReq(ctx, conn, m, uid, platformID)
	default:
	}
}

func (l *MsggatewayLogic) sendErrMsg(ctx context.Context, conn *UserConn, code int32, errMsg string, uid string, platformID string) {
	mReply := &pb.BodyResp{
		ReqIdentifier: types.WSErrorMsg,
		ErrCode:       uint32(code),
		ErrMsg:        errMsg,
	}
	l.sendMsg(ctx, conn, mReply, uid, platformID)
}

func (l *MsggatewayLogic) sendMsg(ctx context.Context, conn *UserConn, resp *pb.BodyResp, uid string, platformID string) {
	b, err := proto.Marshal(resp)
	if err != nil {
		logx.WithContext(ctx).Error(resp.ReqIdentifier, " ", resp.ErrCode, " ", resp.ErrMsg, " Encode Msg error ", conn.RemoteAddr().String(), " ", uid, " ", platformID, " ", err.Error())
		return
	}
	err = l.writeMsg(conn, websocket.BinaryMessage, b)
	if err != nil {
		logx.WithContext(ctx).Error(resp.ReqIdentifier, " ", resp.ErrCode, " ", resp.ErrMsg, " WS WriteMsg error ", conn.RemoteAddr().String(), " ", uid, " ", platformID, " ", err.Error())
	}
}

func (l *MsggatewayLogic) writeMsg(conn *UserConn, a int, msg []byte) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	return conn.WriteMessage(a, msg)
}

func (l *MsggatewayLogic) sendMsgReq(ctx context.Context, conn *UserConn, m *pb.BodyReq, uid string, platformID string) {
	sendMsgAllCount++
	logx.WithContext(ctx).Info("Ws call success to sendMsgReq start", m.ReqIdentifier, m.SendID, m.Data)
	nReply := new(chatpb.SendMsgResp)
	isPass, errCode, errMsg, pData := l.argsValidate(m, types.WSSendMsg)
	if isPass {
		pbData := pData.(chatpb.SendMsgReq)
		// 是否开启限流
		if l.svcCtx.Config.SendMsgRateLimit.Enable {
			if !l.sendMsgRateLimit(ctx, conn, m, uid, platformID, &pbData) {
				return
			}
		}
		logx.WithContext(ctx).Info("Ws call success to sendMsgReq middle", m.ReqIdentifier, m.SendID, pbData.String())

		reply, err := l.svcCtx.MsgRpc().SendMsg(ctx, &pbData)
		if err != nil {
			logx.WithContext(ctx).Error("UserSendMsg err ", err.Error())
			nReply.ErrCode = types.ErrCodeFailed
			nReply.ErrMsg = err.Error()
			nReply.ClientMsgID = pbData.MsgData.ClientMsgID
			nReply.ReceiveID = pbData.MsgData.ReceiveID
			nReply.ConversationType = pbData.MsgData.ConversationType
			nReply.ContentType = pbData.MsgData.ContentType
			l.sendMsgResp(ctx, conn, m, nReply, uid, platformID)
		} else {
			logx.WithContext(ctx).Info("rpc call success to sendMsgReq", reply.String())
			l.sendMsgResp(ctx, conn, m, reply, uid, platformID)
		}
	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		l.sendMsgResp(ctx, conn, m, nReply, uid, platformID)
	}
}

func (l *MsggatewayLogic) sendMsgResp(ctx context.Context, conn *UserConn, m *pb.BodyReq, resp *chat.SendMsgResp, uid string, platformID string) {
	b, _ := proto.Marshal(resp)
	mReply := &pb.BodyResp{
		ReqIdentifier: m.ReqIdentifier,
		ErrCode:       uint32(resp.ErrCode),
		ErrMsg:        resp.ErrMsg,
		Data:          b,
	}
	l.sendMsg(ctx, conn, mReply, uid, platformID)
}

func (l *MsggatewayLogic) pullMsgBySeqListReq(ctx context.Context, conn *UserConn, m *pb.BodyReq, uid string, platformID string) {
	logx.WithContext(ctx).Info("Ws call success to pullMsgBySeqListReq start", m.SendID, m.ReqIdentifier, m.Data)
	nReply := new(chatpb.PullMsgListResp)
	isPass, errCode, errMsg, data := l.argsValidate(m, types.WSPullMsgBySeqList)
	if isPass {
		rpcReq := data.(chatpb.PullMsgBySeqListReq)
		rpcReq.UserID = uid
		logx.WithContext(ctx).Info("Ws call success to pullMsgBySeqListReq middle", m.SendID, m.ReqIdentifier, data.(chatpb.PullMsgBySeqListReq).SeqList)
		reply, err := l.svcCtx.MsgRpc().PullMessageBySeqList(ctx, &rpcReq)
		if err != nil {
			logx.WithContext(ctx).Errorf("pullMsgBySeqListReq err", err.Error())
			nReply.ErrCode = types.ErrCodeFailed
			nReply.ErrMsg = err.Error()
			l.pullMsgBySeqListResp(ctx, conn, m, nReply, uid, platformID)
		} else {
			logx.WithContext(ctx).Info("rpc call success to pullMsgBySeqListReq ", reply.String())
			l.pullMsgBySeqListResp(ctx, conn, m, reply, uid, platformID)
		}
	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		l.pullMsgBySeqListResp(ctx, conn, m, nReply, uid, platformID)
	}
}

func (l *MsggatewayLogic) pullMsgByGroupSeqListReq(ctx context.Context, conn *UserConn, m *pb.BodyReq, uid string, platformID string) {
	logx.WithContext(ctx).Info("Ws call success to pullMsgByGroupSeqListReq start", m.SendID, m.ReqIdentifier, m.Data)
	nReply := new(chatpb.PullMsgListResp)
	isPass, errCode, errMsg, data := l.argsValidate(m, types.WSPullMsgByGroupSeqList)
	if isPass {
		rpcReq := data.(chatpb.PullMsgByGroupSeqListReq)
		logx.WithContext(ctx).Info("Ws call success to pullMsgBySeqListReq middle", m.SendID, m.ReqIdentifier, rpcReq.SeqList)
		reply, err := l.svcCtx.MsgRpc().PullMessageByGroupSeqList(ctx, &rpcReq)
		if err != nil {
			logx.WithContext(ctx).Errorf("pullMsgBySeqListReq err", err.Error())
			nReply.ErrCode = types.ErrCodeFailed
			nReply.ErrMsg = err.Error()
			l.pullMsgByGroupSeqListResp(ctx, conn, m, nReply, uid, platformID)
		} else {
			logx.WithContext(ctx).Info("rpc call success to pullMsgBySeqListReq ", reply.String())
			l.pullMsgByGroupSeqListResp(ctx, conn, m, reply, uid, platformID)
		}
	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		l.pullMsgByGroupSeqListResp(ctx, conn, m, nReply, uid, platformID)
	}
}

func (l *MsggatewayLogic) pullMsgBySeqListResp(ctx context.Context, conn *UserConn, m *pb.BodyReq, resp *chatpb.PullMsgListResp, uid string, platformID string) {
	logx.WithContext(ctx).Info("pullMsgBySeqListResp come  here ", resp.String())
	c, _ := proto.Marshal(resp)
	mReply := &pb.BodyResp{
		ReqIdentifier: m.ReqIdentifier,
		ErrCode:       uint32(resp.ErrCode),
		ErrMsg:        resp.ErrMsg,
		Data:          c,
	}
	logx.WithContext(ctx).Info("pullMsgBySeqListResp all data  is ", mReply.ReqIdentifier, mReply.ErrCode, mReply.ErrMsg,
		len(mReply.Data))

	l.sendMsg(ctx, conn, mReply, uid, platformID)
}

func (l *MsggatewayLogic) pullMsgByGroupSeqListResp(ctx context.Context, conn *UserConn, m *pb.BodyReq, resp *chatpb.PullMsgListResp, uid string, platformID string) {
	logx.WithContext(ctx).Info("pullMsgByGroupSeqListResp come  here ", resp.String())
	for _, data := range resp.List {
		logx.Infof("pullMsgByGroupSeqListResp data is %s", data.String())
	}
	c, _ := proto.Marshal(resp)
	mReply := &pb.BodyResp{
		ReqIdentifier: m.ReqIdentifier,
		ErrCode:       uint32(resp.ErrCode),
		ErrMsg:        resp.ErrMsg,
		Data:          c,
	}
	logx.WithContext(ctx).Info("pullMsgBySeqListResp all data  is ", mReply.ReqIdentifier, mReply.ErrCode, mReply.ErrMsg,
		len(mReply.Data))

	l.sendMsg(ctx, conn, mReply, uid, platformID)
}
