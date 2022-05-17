package wslogic

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/showurl/Zero-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

func (l *MsggatewayLogic) addUserConn(uid string, platformID string, conn *UserConn, token string) {
	rwLock.Lock()
	defer rwLock.Unlock()
	if oldConnMap, ok := l.wsUserToConn[uid]; ok {
		oldConnMap[platformID] = conn
		l.wsUserToConn[uid] = oldConnMap
	} else {
		i := make(map[string]*UserConn)
		i[platformID] = conn
		l.wsUserToConn[uid] = i
	}
	if oldStringMap, ok := l.wsConnToUser[conn]; ok {
		oldStringMap[platformID] = uid
		l.wsConnToUser[conn] = oldStringMap
	} else {
		i := make(map[string]string)
		i[platformID] = uid
		l.wsConnToUser[conn] = i
	}
	count := 0
	for _, v := range l.wsUserToConn {
		count = count + len(v)
	}
}

func (l *MsggatewayLogic) getUserUid(conn *UserConn) (uid string, platform string) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if stringMap, ok := l.wsConnToUser[conn]; ok {
		for k, v := range stringMap {
			platform = k
			uid = v
		}
		return uid, platform
	}
	return "", ""
}

func (l *MsggatewayLogic) delUserConn(conn *UserConn) {
	rwLock.Lock()
	defer rwLock.Unlock()
	var platform, uid string
	if oldStringMap, ok := l.wsConnToUser[conn]; ok {
		for k, v := range oldStringMap {
			platform = k
			uid = v
		}
		if oldConnMap, ok := l.wsUserToConn[uid]; ok {
			delete(oldConnMap, platform)
			l.wsUserToConn[uid] = oldConnMap
			if len(oldConnMap) == 0 {
				delete(l.wsUserToConn, uid)
			}
			count := 0
			for _, v := range l.wsUserToConn {
				count = count + len(v)
			}
		}
		delete(l.wsConnToUser, conn)
	}
	err := conn.Close()
	if err != nil {
		logx.WithContext(l.ctx).Error("close conn err", "", "uid", uid, "platform", platform)
	}
}

func (l *MsggatewayLogic) GetUserConn(uid string, platform string) *UserConn {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if connMap, ok := l.wsUserToConn[uid]; ok {
		if conn, flag := connMap[platform]; flag {
			return conn
		}
	}
	return nil
}

func (l *MsggatewayLogic) SendMsgToUser(ctx context.Context, conn *UserConn, bMsg []byte, in *pb.OnlinePushMsgReq, RecvPlatForm, RecvID string) (ResultCode int64) {
	err := l.writeMsg(conn, websocket.BinaryMessage, bMsg)
	if err != nil {
		logx.WithContext(ctx).Error("send msg to user err ", "", "err ", err.Error())
		ResultCode = -2
		return ResultCode
	} else {
		ResultCode = 0
		return ResultCode
	}
}
