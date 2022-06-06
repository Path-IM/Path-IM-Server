package wslogic

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wsrepository"
	msggatewaypb "github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"sync"
	"time"

	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/types"
	"github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/internal/wssvc"
	commonTypes "github.com/Path-IM/Path-IM-Server/common/types"

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
	rep          *wsrepository.Rep
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
		rep:    wsrepository.NewRep(svcCtx),
	}
	ws.wsMaxConnNum = ws.svcCtx.Config.Websocket.MaxConnNum
	ws.wsConnToUser = make(map[*UserConn]map[string]string)
	ws.wsUserToConn = make(map[string]map[string]*UserConn)
	ws.wsUpGrader = &websocket.Upgrader{
		HandshakeTimeout: time.Duration(ws.svcCtx.Config.Websocket.TimeOut) * time.Second,
		ReadBufferSize:   ws.svcCtx.Config.Websocket.ReadBufferSize,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
	msgGatewayLogic = ws
	return msgGatewayLogic
}

func (l *MsggatewayLogic) Msggateway(req *types.Request) (*types.Response, bool) {
	if len(req.Token) != 0 && len(req.UserID) != 0 && len(req.Platform) != 0 {
		// 调用rpc验证token
		resp, err := l.svcCtx.ImUserService().VerifyToken(l.ctx, &pb.VerifyTokenReq{
			Token:    req.Token,
			Platform: req.Platform,
			SendID:   req.UserID,
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
	l.addUserConn(uid, req.Platform, newConn, req.Token)
	go l.readMsg(newConn, uid, req.Platform)
	return nil
}

func (l *MsggatewayLogic) readMsg(conn *UserConn, uid string, platformID string) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if messageType == websocket.PingMessage {
			l.sendMsg(context.Background(), conn, &msggatewaypb.BodyResp{
				ReqIdentifier: 0,
				ErrCode:       0,
				ErrMsg:        "",
				Data:          []byte("pong"),
			}, uid, platformID)
		}
		if err != nil {
			uid, platform := l.getUserUid(conn)
			logx.Error("WS ReadMsg error ", " userIP ", conn.RemoteAddr().String(), " userUid ", uid, " platform ", platform, " error ", err.Error())
			l.delUserConn(conn)
			return
		}
		xtrace.RunWithTrace("", func(ctx context.Context) {
			l.msgParse(ctx, conn, msg, uid, platformID)
		}, attribute.KeyValue{
			Key:   "uid",
			Value: attribute.StringValue(uid),
		}, attribute.KeyValue{
			Key:   "platformID",
			Value: attribute.StringValue(platformID),
		})
	}
}

func (l *MsggatewayLogic) getSeqReq(ctx context.Context, conn *UserConn, m *msggatewaypb.BodyReq, uid string, platformID string) {
	rpcReq := chatpb.GetMinAndMaxSeqReq{}
	nReply := new(chatpb.GetMinAndMaxSeqResp)
	rpcReq.UserID = uid
	rpcReply, err := l.svcCtx.MsgRpc().GetMaxAndMinSeq(ctx, &rpcReq)
	if err != nil {
		logx.WithContext(ctx).Error("rpc call failed to getSeqReq", err, rpcReq.String())
		nReply.ErrCode = commonTypes.ErrCodeFailed
		nReply.ErrMsg = err.Error()
		l.getSeqResp(ctx, conn, m, nReply, uid, platformID)
	} else {
		logx.WithContext(ctx).Info("rpc call success to getSeqReq", rpcReply.String())
		l.getSeqResp(ctx, conn, m, rpcReply, uid, platformID)
	}
}
func (l *MsggatewayLogic) getGroupSeqReq(ctx context.Context, conn *UserConn, m *msggatewaypb.BodyReq, uid string, platformID string) {
	rpcReq := &chatpb.GetMinAndMaxGroupSeqReq{}
	err := proto.Unmarshal(m.Data, rpcReq)
	nReply := new(chatpb.GetMinAndMaxGroupSeqResp)
	if err != nil {
		logx.WithContext(ctx).Error("proto.Unmarshal failed ", err)
		nReply.ErrCode = commonTypes.ErrCodeParams
		nReply.ErrMsg = "param verify failed"
		l.getGroupResp(ctx, conn, m, nReply, uid, platformID)
	}
	rpcReply, err := l.svcCtx.MsgRpc().GetMinAndMaxGroupSeq(ctx, rpcReq)
	if err != nil {
		logx.WithContext(ctx).Error("rpc call failed to getSeqReq", err, rpcReq.String())
		nReply.ErrCode = commonTypes.ErrCodeFailed
		nReply.ErrMsg = err.Error()
		l.getGroupResp(ctx, conn, m, nReply, uid, platformID)
	} else {
		logx.WithContext(ctx).Info("rpc call success to getSeqReq", rpcReply.String())
		l.getGroupResp(ctx, conn, m, rpcReply, uid, platformID)
	}
}

func (l *MsggatewayLogic) getSeqResp(ctx context.Context, conn *UserConn, m *msggatewaypb.BodyReq, resp *chatpb.GetMinAndMaxSeqResp, uid string, platformID string) {
	b, _ := proto.Marshal(resp)
	mReply := &msggatewaypb.BodyResp{
		ReqIdentifier: m.ReqIdentifier,
		ErrCode:       resp.ErrCode,
		ErrMsg:        resp.ErrMsg,
		Data:          b,
	}
	l.sendMsg(ctx, conn, mReply, uid, platformID)
}

func (l *MsggatewayLogic) getGroupResp(ctx context.Context, conn *UserConn, m *msggatewaypb.BodyReq, resp *chatpb.GetMinAndMaxGroupSeqResp, uid string, platformID string) {
	b, _ := proto.Marshal(resp)
	mReply := &msggatewaypb.BodyResp{
		ReqIdentifier: m.ReqIdentifier,
		ErrCode:       resp.ErrCode,
		ErrMsg:        resp.ErrMsg,
		Data:          b,
	}
	l.sendMsg(ctx, conn, mReply, uid, platformID)
}
