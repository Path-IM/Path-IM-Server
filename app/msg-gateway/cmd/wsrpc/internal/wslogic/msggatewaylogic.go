package wslogic

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/showurl/Zero-IM-Server/app/im-user/cmd/rpc/pb"
	msggatewaypb "github.com/showurl/Zero-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	chatpb "github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/showurl/Zero-IM-Server/common/xtrace"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"sync"
	"time"

	"github.com/showurl/Zero-IM-Server/app/msg-gateway/cmd/wsrpc/internal/types"
	"github.com/showurl/Zero-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wssvc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserConn struct {
	*websocket.Conn
	w *sync.Mutex
}

type MsggatewayLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *wssvc.ServiceContext
	wsMaxConnNum int
	wsUpGrader   *websocket.Upgrader
	wsConnToUser map[*UserConn]map[string]string
	wsUserToConn map[string]map[string]*UserConn
}

var msgGatewayLogic *MsggatewayLogic

func NewMsggatewayLogic(ctx context.Context, svcCtx *wssvc.ServiceContext) *MsggatewayLogic {
	if msgGatewayLogic != nil {
		return msgGatewayLogic
	}
	ws := &MsggatewayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
	ws.wsMaxConnNum = ws.svcCtx.Config.Websocket.MaxConnNum
	ws.wsConnToUser = make(map[*UserConn]map[string]string)
	ws.wsUserToConn = make(map[string]map[string]*UserConn)
	ws.wsUpGrader = &websocket.Upgrader{
		HandshakeTimeout: time.Duration(ws.svcCtx.Config.Websocket.TimeOut) * time.Second,
		ReadBufferSize:   ws.svcCtx.Config.Websocket.MaxMsgLen,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
	msgGatewayLogic = ws
	return msgGatewayLogic
}

func (l *MsggatewayLogic) Msggateway(req *types.Request) (*types.Response, bool) {
	if len(req.Token) != 0 && len(req.SendID) != 0 && len(req.PlatformID) != 0 {
		// 调用rpc验证token
		resp, err := l.svcCtx.ImUserService.VerifyToken(l.ctx, &pb.VerifyTokenReq{
			Token:    req.Token,
			Platform: req.PlatformID,
			SendID:   req.SendID,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("调用 VerifyToken 失败, err: %s", err.Error())
			return &types.Response{
				Uid:     "",
				ErrMsg:  "调用 VerifyToken 失败",
				Success: false,
			}, false
		}
		if !resp.Success {
			logx.WithContext(l.ctx).Infof("VerifyToken 失败, err: %s", resp.ErrMsg)
			return &types.Response{
				Uid:     resp.Uid,
				ErrMsg:  resp.ErrMsg,
				Success: false,
			}, false
		}
		return &types.Response{
			Uid:     resp.Uid,
			ErrMsg:  "",
			Success: true,
		}, true
	}
	return &types.Response{
		Uid:     "",
		ErrMsg:  "参数错误",
		Success: false,
	}, false
}

func (l *MsggatewayLogic) WsUpgrade(uid string, req *types.Request, w http.ResponseWriter, r *http.Request, header http.Header) error {
	conn, err := l.wsUpGrader.Upgrade(w, r, header)
	if err != nil {
		return err
	}
	newConn := &UserConn{conn, new(sync.Mutex)}
	userCount++
	l.addUserConn(uid, req.PlatformID, newConn, req.Token)
	go l.readMsg(newConn, uid, req.PlatformID)
	return nil
}

func (l *MsggatewayLogic) readMsg(conn *UserConn, uid string, platformID string) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if messageType == websocket.PingMessage {
			l.sendMsg(context.Background(), conn, Resp{
				ReqIdentifier: 0,
				MsgIncr:       "",
				ErrCode:       0,
				ErrMsg:        "",
				Data:          []byte("pong"),
			})
		}
		if err != nil {
			uid, platform := l.getUserUid(conn)
			logx.Error("WS ReadMsg error ", " userIP ", conn.RemoteAddr().String(), " userUid ", uid, " platform ", platform, " error ", err.Error())
			userCount--
			l.delUserConn(conn)
			return
		}
		xtrace.RunWithTrace("", func(ctx context.Context) {
			l.msgParse(ctx, conn, msg)
		}, attribute.KeyValue{
			Key:   "uid",
			Value: attribute.StringValue(uid),
		}, attribute.KeyValue{
			Key:   "platformID",
			Value: attribute.StringValue(platformID),
		})
	}
}

func (l *MsggatewayLogic) getSeqReq(ctx context.Context, conn *UserConn, m *msggatewaypb.Req) {
	rpcReq := chatpb.GetMaxAndMinSeqReq{}
	nReply := new(chatpb.GetMaxAndMinSeqResp)
	rpcReq.UserID = m.SendID
	rpcReply, err := l.svcCtx.MsgRpc.GetMaxAndMinSeq(ctx, &rpcReq)
	if err != nil {
		logx.WithContext(ctx).Error("rpc call failed to getSeqReq", err, rpcReq.String())
		nReply.ErrCode = 500
		nReply.ErrMsg = err.Error()
		l.getSeqResp(ctx, conn, m, nReply)
	} else {
		logx.WithContext(ctx).Info("rpc call success to getSeqReq", rpcReply.String())
		l.getSeqResp(ctx, conn, m, rpcReply)
	}
}
func (l *MsggatewayLogic) getSuperGroupSeqReq(ctx context.Context, conn *UserConn, m *msggatewaypb.Req) {
	rpcReq := &chatpb.GetMaxAndMinSuperGroupSeqReq{}
	err := proto.Unmarshal(m.Data, rpcReq)
	nReply := new(chatpb.GetMaxAndMinSuperGroupSeqResp)
	if err != nil {
		logx.WithContext(ctx).Error("proto.Unmarshal failed ", err)
		nReply.ErrCode = 300
		nReply.ErrMsg = "param verify failed"
		l.getSuperGroupResp(ctx, conn, m, nReply)
	}
	rpcReply, err := l.svcCtx.MsgRpc.GetSuperGroupMaxAndMinSeq(ctx, rpcReq)
	if err != nil {
		logx.WithContext(ctx).Error("rpc call failed to getSeqReq", err, rpcReq.String())
		nReply.ErrCode = 500
		nReply.ErrMsg = err.Error()
		l.getSuperGroupResp(ctx, conn, m, nReply)
	} else {
		logx.WithContext(ctx).Info("rpc call success to getSeqReq", rpcReply.String())
		l.getSuperGroupResp(ctx, conn, m, rpcReply)
	}
}

func (l *MsggatewayLogic) getSeqResp(ctx context.Context, conn *UserConn, m *msggatewaypb.Req, pb *chatpb.GetMaxAndMinSeqResp) {
	var mReplyData chatpb.GetMaxAndMinSeqResp
	mReplyData.MaxSeq = pb.GetMaxSeq()
	mReplyData.MinSeq = pb.GetMinSeq()
	b, _ := proto.Marshal(&mReplyData)
	mReply := Resp{
		ReqIdentifier: int32(m.ReqIdentifier),
		MsgIncr:       m.MsgIncr,
		ErrCode:       pb.GetErrCode(),
		ErrMsg:        pb.GetErrMsg(),
		Data:          b,
	}
	l.sendMsg(ctx, conn, mReply)
}

func (l *MsggatewayLogic) getSuperGroupResp(ctx context.Context, conn *UserConn, m *msggatewaypb.Req, pb *chatpb.GetMaxAndMinSuperGroupSeqResp) {
	var mReplyData chatpb.GetMaxAndMinSuperGroupSeqResp
	mReplyData.SuperGroupSeqList = pb.GetSuperGroupSeqList()
	b, _ := proto.Marshal(&mReplyData)
	mReply := Resp{
		ReqIdentifier: int32(m.ReqIdentifier),
		MsgIncr:       m.MsgIncr,
		ErrCode:       pb.GetErrCode(),
		ErrMsg:        pb.GetErrMsg(),
		Data:          b,
	}
	l.sendMsg(ctx, conn, mReply)
}
