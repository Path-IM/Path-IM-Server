package wslogic

import (
	"context"
	imuserpb "github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"
	msggatewaypb "github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func (l *MsggatewayLogic) GetUsersOnlineStatus(ctx context.Context, req *msggatewaypb.GetUsersOnlineStatusReq) (*msggatewaypb.GetUsersOnlineStatusResp, error) {
	var resp msggatewaypb.GetUsersOnlineStatusResp
	for _, userID := range req.UserIDList {
		platformMap, err := l.rep.GetUserConn(ctx, userID)
		if err != nil {
			return nil, err
		}
		resp.StatusList = append(resp.StatusList, &msggatewaypb.GetUsersOnlineStatusResp_UserStatus{
			UserID:          userID,
			PlatformAddrMap: platformMap,
		})
	}
	return &resp, nil
}

func (l *MsggatewayLogic) addUserConn(uid string, platformID string, conn *UserConn, token string) error {
	rwLock.Lock()
	defer rwLock.Unlock()
	err := l.rep.AddUserConn(uid, platformID)
	if err != nil {
		return err
	}
	if l.svcCtx.Config.EnableUserCallback {
		_, err = l.svcCtx.ImUserService().UserCallback(l.ctx, &imuserpb.UserCallbackReq{
			Event:      imuserpb.UserCallbackReq_Online,
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

func (l *MsggatewayLogic) delUserConn(ctx context.Context, conn *UserConn) {
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
	err = l.rep.DelUserConn(ctx, uid, platform)
	if err != nil {
		logx.WithContext(l.ctx).Error("redis DelUserConn err ", err.Error(), " uid ", uid, " platform ", platform)
	}
	if l.svcCtx.Config.EnableUserCallback {
		_, err = l.svcCtx.ImUserService().UserCallback(l.ctx, &imuserpb.UserCallbackReq{
			Event:      imuserpb.UserCallbackReq_Offline,
			Timestamp:  time.Now().UnixMilli(),
			UserID:     uid,
			Platform:   platform,
			RemoteAddr: conn.RemoteAddr().String(),
		})
		if err != nil {
			logx.Errorf("user callback err %s", err.Error())
			err = nil
		}
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

func (l *MsggatewayLogic) DelUserConn(ctx context.Context, uid string, platform string) {
	conn := l.GetUserConn(uid, platform)
	if conn != nil {
		l.delUserConn(ctx, conn)
	}
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
