package logic

import (
	msgcallbackpb "github.com/showurl/Zero-IM-Server/app/msg-callback/cmd/rpc/pb"
	chatpb "github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/showurl/Zero-IM-Server/common/types"
)

func (l *SendMsgLogic) copyCallbackCommonReqStruct(msg *chatpb.SendMsgReq) msgcallbackpb.CommonCallbackReq {
	return msgcallbackpb.CommonCallbackReq{
		SendID:           msg.MsgData.SendID,
		ServerMsgID:      msg.MsgData.ServerMsgID,
		ClientMsgID:      msg.MsgData.ClientMsgID,
		SenderPlatformID: msg.MsgData.SenderPlatformID,
		SenderNickname:   msg.MsgData.SenderNickname,
		SessionType:      msg.MsgData.SessionType,
		MsgFrom:          msg.MsgData.MsgFrom,
		ContentType:      msg.MsgData.ContentType,
		Status:           msg.MsgData.Status,
		CreateTime:       msg.MsgData.CreateTime,
		Content:          string(msg.MsgData.Content),
	}
}

func (l *SendMsgLogic) callbackWordFilter(msg *chatpb.SendMsgReq) (canSend bool, err error) {
	if !l.svcCtx.Config.Callback.CallbackWordFilter.Enable || msg.MsgData.ContentType != types.Text {
		return true, nil
	}
	commonCallbackReq := l.copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = msgcallbackpb.CallbackCommand_WordFilter
	req := msgcallbackpb.CallbackWordFilterReq{
		CommonCallbackReq: &commonCallbackReq,
	}
	resp, err := l.svcCtx.MsgCallback.CallbackWordFilter(l.ctx, &req)
	if err != nil {
		if l.svcCtx.Config.Callback.CallbackWordFilter.ContinueOnError {
			return true, err
		} else {
			return false, err
		}
	}
	if resp.CommonCallbackResp.ActionCode == msgcallbackpb.ActionCode_Forbidden && resp.CommonCallbackResp.ErrCode == msgcallbackpb.ErrCode_HandleSuccess {
		return false, nil
	}
	if resp.CommonCallbackResp.ErrCode == msgcallbackpb.ErrCode_HandleSuccess {
		msg.MsgData.Content = []byte(resp.ReplaceContent)
	}
	return true, err
}
