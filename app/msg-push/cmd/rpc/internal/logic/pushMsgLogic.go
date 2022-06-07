package logic

import (
	"context"
	"encoding/json"
	imuserpb "github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/internal/repository"
	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/internal/svc"
	pushpb "github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type PushMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	rep *repository.Rep
}

func NewPushMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushMsgLogic {
	return &PushMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		rep:    repository.NewRep(svcCtx),
	}
}

func (l *PushMsgLogic) OfflinePushMsg(in *pushpb.PushMsgReq) (resp *pushpb.PushMsgResp, err error) {
	// 判断是单聊还是群聊
	if in.MsgData.ConversationType == types.SingleChatType {
		uid := in.PushToUserID[0]
		if in.MsgData.SendID == uid {
			return nil, err
		}
		optResp, err := l.svcCtx.ImUserService.GetSingleConversationRecvMsgOpts(
			l.ctx,
			&imuserpb.GetSingleConversationRecvMsgOptsReq{
				UserID:       uid,
				SenderUserID: in.MsgData.SendID,
			},
		)
		if err != nil {
			return nil, err
		}
		if optResp.Opts == imuserpb.RecvMsgOpt_ReceiveMessage {
			// 是否在线
			online, err := l.rep.IsUserOnline(uid)
			if err != nil {
				l.Logger.Errorf("redis IsUserOnline err: %v", err)
				err = nil
				online = true
			}
			if !online {
				// 离线推送
				// 用户预览消息选项
				title := l.svcCtx.Config.OfflinePushDefaultTitle
				buf, _ := json.Marshal(in.MsgData.OfflinePush)
				previewResp, err := l.svcCtx.ImUserService.IfPreviewMessage(l.ctx, &imuserpb.IfPreviewMessageReq{
					SenderID:   in.MsgData.SendID,
					ReceiverID: uid,
					GroupID:    "",
				})
				if err != nil {
					l.Logger.Errorf("IfPreviewMessage err: %v", err)
					err = nil
				} else {
					if previewResp.IfPreview {
						title = in.MsgData.OfflinePush.Title
					} else {
						title = previewResp.ReplaceTitle
					}
				}
				res, e := l.svcCtx.GetOfflinePusher().Push(l.ctx, []string{uid}, title, string(buf))
				if e != nil {
					l.Logger.Errorf("GetOfflinePusher.Push err: %v, %s", e, res)
				} else {
					l.Logger.Infof("GetOfflinePusher.Push res: %s", res)
				}
			}
		}
	} else if in.MsgData.ConversationType == types.GroupChatType {
		err = l.offlinePushGroupMsg(in)
		if err != nil {
			l.Errorf("offlinePushGroupMsg err : %v", err)
		}
	}
	return nil, err
}

func (l *PushMsgLogic) offlinePushGroupMsg(in *pushpb.PushMsgReq) error {
	var onlineUserId []string
	var offlineUserId []string
	for _, pushUser := range in.PushToUserID {
		if online, err := l.rep.IsUserOnline(pushUser); err != nil || online {
			onlineUserId = append(onlineUserId, pushUser)
		} else {
			offlineUserId = append(offlineUserId, pushUser)
		}
	}
	// 离线推送
	// 用户预览消息选项
	title := l.svcCtx.Config.OfflinePushGroupTitle
	buf, _ := json.Marshal(in.MsgData.OfflinePush)
	res, e := l.svcCtx.GetOfflinePusher().Push(l.ctx, offlineUserId, title, string(buf))
	if e != nil {
		l.Logger.Errorf("GetOfflinePusher.Push err: %v, %s", e, res)
	} else {
		l.Logger.Infof("GetOfflinePusher.Push res: %s", res)
	}
	return e
}
