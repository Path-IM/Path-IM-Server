package wslogic

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

func (l *MsggatewayLogic) addUserConn(uid string, platformID string, conn *UserConn, token string) error {
	rwLock.Lock()
	defer rwLock.Unlock()
	err := l.rep.AddUserConn(uid, platformID)
	if err != nil {
		return err
	}
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
	return nil
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
		}
		delete(l.wsConnToUser, conn)
	}
	err := conn.Close()
	if err != nil {
		logx.WithContext(l.ctx).Error("close conn err", "", "uid", uid, "platform", platform)
	}
	err = l.rep.DelUserConn(uid, platform)
	if err != nil {
		logx.WithContext(l.ctx).Error("redis DelUserConn err ", err.Error(), " uid ", uid, " platform ", platform)
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
func (l *MsggatewayLogic) GetOnlineUserMap() map[string]interface{} {
	m := make(map[string]interface{})
	rwLock.RLock()
	defer rwLock.RUnlock()
	for uid := range l.wsUserToConn {
		m[uid] = nil
	}
	return m
}

func (l *MsggatewayLogic) DelUserConn(uid string, platform string) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	conn := l.GetUserConn(uid, platform)
	l.delUserConn(conn)
}

func (l *MsggatewayLogic) SendMsgToUser(ctx context.Context, conn *UserConn, bMsg []byte, RecvPlatForm, RecvID string) error {
	//logx.WithContext(ctx).Infof("发送给用户:%s[%s],%+v", RecvID, RecvPlatForm, string(bMsg))
	err := l.writeMsg(conn, websocket.BinaryMessage, bMsg)
	if err != nil {
		logx.WithContext(ctx).Error("send msg to user err ", "", "err ", err.Error())
		return err
	} else {
		return nil
	}
}
