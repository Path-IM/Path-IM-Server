package logic

import (
	"context"
	"encoding/json"
	imuserpb "github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/pb"
	onlinemessagerelayservice "github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/onlineMessageRelayService"
	gatewaypb "github.com/Path-IM/Path-IM-Server/app/msg-gateway/cmd/wsrpc/pb"
	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/pb"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/types"
	"github.com/Path-IM/Path-IM-Server/common/utils"
	numUtils "github.com/Path-IM/Path-IM-Server/common/utils/num"
	strUtils "github.com/Path-IM/Path-IM-Server/common/utils/str"
	"github.com/Path-IM/Path-IM-Server/common/xtrace"
	"github.com/zeromicro/go-zero/core/mr"
	"go.opentelemetry.io/otel/attribute"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	successCount = uint64(0)
	pushTerminal = []int32{types.IOSPlatformID, types.AndroidPlatformID}
)

type PushMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushMsgLogic {
	return &PushMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushMsgLogic) PushMsg(in *pb.PushMsgReq) (*pb.PushMsgResp, error) {
	l.MsgToUser(in)
	return &pb.PushMsgResp{ResultCode: 0}, nil
}

func (l *PushMsgLogic) getAllMsgGatewayService() (services []onlinemessagerelayservice.OnlineMessageRelayService, err error) {
	if l.svcCtx.Config.MsgGatewayRpcEtcd != nil {
		return onlinemessagerelayservice.GetAllByEtcd(l.ctx, *l.svcCtx.Config.MsgGatewayRpcEtcd, l.svcCtx.Config.MsgGatewayRpcEtcd.Key)
	} else if l.svcCtx.Config.MsgGatewayRpcK8sTarget != "" {
		return onlinemessagerelayservice.GetAllByK8s(l.svcCtx.Config.MsgGatewayRpcK8sTarget)
	} else {
		return onlinemessagerelayservice.GetAllByEndpoints(l.ctx, l.svcCtx.Config.MsgGatewayRpcEndpoints)
	}
}

func (l *PushMsgLogic) MsgToUser(pushMsg *pb.PushMsgReq) {
	var wsResult []*gatewaypb.SingleMsgToUser
	isOfflinePush := utils.GetSwitchFromOptions(pushMsg.MsgData.MsgOptions, types.IsOfflinePush)

	services, err := l.getAllMsgGatewayService()
	if err != nil {
		l.Errorf("getAllMsgGatewayService error: %v", err)
		err = nil
	}
	var fs []func() error
	for index, msgClient := range services {
		fs = append(fs, func() error {
			var reply *gatewaypb.OnlinePushMsgResp
			var err error
			xtrace.StartFuncSpan(l.ctx, "MsgToUser.OnlinePushMsg", func(ctx context.Context) {
				reply, err = msgClient.OnlinePushMsg(ctx, &gatewaypb.OnlinePushMsgReq{MsgData: pushMsg.MsgData, PushToUserID: pushMsg.PushToUserID})
			}, attribute.Int("index", index))
			if err != nil {
				l.Errorf("OnlinePushMsg error: %v", err)
				return nil
			}
			if reply != nil && reply.Resp != nil {
				wsResult = append(wsResult, reply.Resp...)
			}
			return nil
		})
	}
	_ = mr.Finish(fs...)
	l.Info("push_result ", wsResult, " sendData ", pushMsg.MsgData)
	successCount++
	if isOfflinePush && pushMsg.PushToUserID != pushMsg.MsgData.SendID {
		for _, v := range wsResult {
			if v.ResultCode == 0 {
				continue
			}
			if numUtils.IsContainInt32(v.RecvPlatFormID, pushTerminal) {
				var UIDList []string
				UIDList = append(UIDList, v.RecvID)
				bCustomContent, _ := json.Marshal(pushMsg.MsgData.OfflinePush)
				jsonCustomContent := string(bCustomContent)
				title := pushMsg.MsgData.OfflinePush.Title
				// 判断接受者是否开启了预览消息
				message, err := l.svcCtx.ImUserService.IfPreviewMessage(l.ctx, &imuserpb.IfPreviewMessageReq{
					SenderID:   pushMsg.MsgData.SendID,
					ReceiverID: pushMsg.PushToUserID,
				})
				if err != nil || message.CommonResp == nil || message.CommonResp.ErrCode != 0 {
					l.Errorf("IfPreviewMessage error: %v", err)
					title = l.svcCtx.Config.OfflinePushDefaultTitle
				} else {
					if message.CommonResp.ErrCode == 0 && !message.IfPreview {
						title = message.ReplaceTitle
					}
				}
				xtrace.StartFuncSpan(l.ctx, "MsgToUser.OfflinePushMsg", func(ctx context.Context) {
					_, err = l.svcCtx.GetOfflinePusher().Push(ctx, UIDList, title, jsonCustomContent)
				})
				if err != nil {
					l.Error("offline push error ", pushMsg.String(), err.Error())
				}
				break
			}
		}
	}
}

func (l *PushMsgLogic) PushGroupMsg(in *chatpb.PushMsgDataToMQ) (*pb.PushMsgResp, error) {
	isOfflinePush := utils.GetSwitchFromOptions(in.MsgData.MsgOptions, types.IsOfflinePush)

	tagAll := false
	// 如果艾特人了
	if in.MsgData.AtUserIDList != nil {
		tagAll = strUtils.IsContain(types.AtAllString, in.MsgData.AtUserIDList)
	}
	// 被艾特的人 先去获取被艾特的人是否屏蔽了群消息
	var atUsers *imuserpb.GetUserListFromGroupWithOptResp
	atPushUserChan := make(chan string, 1)
	go l.listenAtPushUserChan(atPushUserChan, in)
	var err error
	if tagAll {
		xtrace.StartFuncSpan(l.ctx, "PushGroupMsg.GetUserListFromGroupWithOpt", func(ctx context.Context) {
			// 我们去查询这个群的所有接收消息通知的用户
			atUsers, err = l.svcCtx.ImUserService.GetUserListFromGroupWithOpt(l.ctx, &imuserpb.GetUserListFromGroupWithOptReq{
				GroupID: in.MsgData.ReceiveID,
				Opts: []imuserpb.RecvMsgOpt{
					imuserpb.RecvMsgOpt_ReceiveMessage,
					imuserpb.RecvMsgOpt_ReceiveNotNotifyMessage,
				},
			})
		})
		if err == nil {
			l.pushGroupMsg(in, atUsers, nil, isOfflinePush, atPushUserChan)
		} else {
			logx.WithContext(l.ctx).Error("GetUserListFromGroupWithOpt failed, err: ", err)
			err = nil
		}
	} else if len(in.MsgData.AtUserIDList) > 0 {
		xtrace.StartFuncSpan(l.ctx, "PushGroupMsg.GetUserListFromGroupWithOpt", func(ctx context.Context) {
			// 我们去查询这个群的所有接收消息通知的用户
			atUsers, err = l.svcCtx.ImUserService.GetUserListFromGroupWithOpt(l.ctx, &imuserpb.GetUserListFromGroupWithOptReq{
				GroupID: in.MsgData.ReceiveID,
				Opts: []imuserpb.RecvMsgOpt{
					imuserpb.RecvMsgOpt_ReceiveNotNotifyMessage,
				},
				UserIDList: in.MsgData.AtUserIDList,
			})
		})
		if err == nil {
			var verifyAtUsers = &imuserpb.GetUserListFromGroupWithOptResp{
				CommonResp:    &imuserpb.CommonResp{},
				UserIDOptList: nil,
			}
			for _, opt := range atUsers.UserIDOptList {
				if strUtils.IsContain(opt.UserID, in.MsgData.AtUserIDList) {
					verifyAtUsers.UserIDOptList = append(verifyAtUsers.UserIDOptList, opt)
				}
			}
			l.pushGroupMsg(in, verifyAtUsers, nil, isOfflinePush, atPushUserChan)
		} else {
			logx.WithContext(l.ctx).Error("GetUserListFromGroupWithOpt failed, err: ", err)
			err = nil
		}
	}
	if tagAll {
		return &pb.PushMsgResp{ResultCode: 0}, nil
	}
	var allUsers *imuserpb.GetUserListFromGroupWithOptResp
	offlinePushUserChan := make(chan string, 1)
	xtrace.StartFuncSpan(l.ctx, "PushGroupMsg.GetUserListFromGroupWithOpt", func(ctx context.Context) {
		// 我们去查询这个群的所有接收消息通知的用户
		allUsers, err = l.svcCtx.ImUserService.GetUserListFromGroupWithOpt(l.ctx, &imuserpb.GetUserListFromGroupWithOptReq{
			GroupID: in.MsgData.ReceiveID,
			Opts: []imuserpb.RecvMsgOpt{
				imuserpb.RecvMsgOpt_ReceiveMessage,
			},
		})
	})
	if err != nil {
		return nil, err
	}
	l.Info("allUsers.UserIDOptList:", allUsers.UserIDOptList)

	go l.listenOfflinePushUserChan(offlinePushUserChan, in)
	l.pushGroupMsg(in, allUsers, in.MsgData.AtUserIDList, isOfflinePush, offlinePushUserChan)
	return &pb.PushMsgResp{ResultCode: 0}, nil
}

func (l *PushMsgLogic) pushGroupMsg(
	in *chatpb.PushMsgDataToMQ,
	users *imuserpb.GetUserListFromGroupWithOptResp,
	atList []string,
	isOfflinePush bool,
	offlinePushUserChan chan string,
) {
	services, _ := l.getAllMsgGatewayService()
	go func() {
		defer func() {
			offlinePushUserChan <- string([]byte{2})
		}()
		for uIndex, user := range users.UserIDOptList {
			if strUtils.IsContain(user.UserID, atList) {
				// 跳过被艾特的人
				continue
			}
			{
				allServiceFailed := true
				var fs []func() error
				for i, service := range services {
					fs = append(fs, func() error {
						allPlatformIsFailed := true
						xtrace.StartFuncSpan(l.ctx, "PushGroupMsg.PushMsgToUser", func(ctx context.Context) {
							resp, err := service.OnlinePushMsg(ctx, &gatewaypb.OnlinePushMsgReq{
								MsgData:      in.MsgData,
								PushToUserID: user.UserID,
							})
							if err != nil {
								l.Errorf("PushMsgToUser error: %v", err)
								return
							}
							if resp == nil || resp.Resp == nil {
								l.Errorf("PushMsgToUser error: resp == nil")
								return
							}
							for _, res := range resp.Resp {
								// 是否全部平台都失败了
								if res.ResultCode != -1 {
									// 成功了
									allPlatformIsFailed = false
									break
								}
							}
						},
							attribute.Int("user.index", uIndex),
							attribute.Int("services.index", i),
							attribute.String("user.id", user.UserID),
						)
						if !allPlatformIsFailed {
							allServiceFailed = false
						}
						return nil
					})
				}
				_ = mr.Finish(fs...)
				if allServiceFailed {
					// 这条消息要不要离线推送
					if isOfflinePush && in.MsgData.SendID != user.UserID {
						offlinePushUserChan <- user.UserID
					}
				}
			}
		}
	}()
}

func (l *PushMsgLogic) listenOfflinePushUserChan(
	userChan chan string,
	pushMsg *chatpb.PushMsgDataToMQ,
) {
	var uids []string
	for uid := range userChan {
		bytes := []byte(uid)
		if len(bytes) == 1 && bytes[0] == 2 {
			break
		}
		uids = append(uids, uid)
	}
	logx.WithContext(l.ctx).Info("开始进行离线推送:", uids)
	bCustomContent, _ := json.Marshal(pushMsg.MsgData.OfflinePush)
	jsonCustomContent := string(bCustomContent)
	title := pushMsg.MsgData.OfflinePush.Title
	// 判断接受者是否开启了预览消息
	message, err := l.svcCtx.ImUserService.IfPreviewMessage(l.ctx, &imuserpb.IfPreviewMessageReq{
		SenderID:   pushMsg.MsgData.SendID,
		ReceiverID: pushMsg.PushToUserID,
		GroupID:    pushMsg.MsgData.ReceiveID,
	})
	if err != nil || message.CommonResp == nil || message.CommonResp.ErrCode != 0 {
		l.Errorf("IfPreviewMessage error: %v", err)
		title = l.svcCtx.Config.OfflinePushDefaultTitle
	} else {
		if message.CommonResp.ErrCode == 0 && !message.IfPreview {
			title = message.ReplaceTitle
		}
	}
	xtrace.StartFuncSpan(l.ctx, "MsgToUser.OfflinePushMsg", func(ctx context.Context) {
		_, err = l.svcCtx.GetOfflinePusher().Push(ctx, uids, title, jsonCustomContent)
	})
	if err != nil {
		l.Error("offline push error ", pushMsg.String(), err.Error())
	}
}

func (l *PushMsgLogic) listenAtPushUserChan(
	userChan chan string,
	pushMsg *chatpb.PushMsgDataToMQ,
) {
	var uids []string
	for uid := range userChan {
		bytes := []byte(uid)
		if len(bytes) == 1 && bytes[0] == 2 {
			break
		}
		uids = append(uids, uid)
	}
	logx.WithContext(l.ctx).Info("开始进行at离线推送:", uids)
	bCustomContent, _ := json.Marshal(pushMsg.MsgData.OfflinePush)
	jsonCustomContent := string(bCustomContent)
	title := pushMsg.MsgData.OfflinePush.Title
	// 判断接受者是否开启了预览消息
	message, err := l.svcCtx.ImUserService.IfPreviewMessage(l.ctx, &imuserpb.IfPreviewMessageReq{
		SenderID:   pushMsg.MsgData.SendID,
		ReceiverID: pushMsg.PushToUserID,
		GroupID:    pushMsg.MsgData.ReceiveID,
	})
	if err != nil || message.CommonResp == nil || message.CommonResp.ErrCode != 0 {
		l.Errorf("IfPreviewMessage error: %v", err)
		title = l.svcCtx.Config.OfflinePushDefaultTitle
	} else {
		if message.CommonResp.ErrCode == 0 && !message.IfPreview {
			title = message.ReplaceTitle
		}
	}
	xtrace.StartFuncSpan(l.ctx, "MsgToUser.OfflinePushMsg", func(ctx context.Context) {
		_, err = l.svcCtx.GetOfflinePusher().Push(ctx, uids, title, jsonCustomContent)
	})
	if err != nil {
		l.Error("offline push error ", pushMsg.String(), err.Error())
	}
}
