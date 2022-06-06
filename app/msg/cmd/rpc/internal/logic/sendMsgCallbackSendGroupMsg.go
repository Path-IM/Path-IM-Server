package logic

import (
	msgcallbackpb "github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/pb"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
)

func (l *SendMsgLogic) copyCallbackCommonReqStruct(msg *chatpb.SendMsgReq) msgcallbackpb.CommonCallbackReq {
	return msgcallbackpb.CommonCallbackReq{
		SendID:           msg.MsgData.SendID,
		ServerMsgID:      msg.MsgData.ServerMsgID,
		ClientMsgID:      msg.MsgData.ClientMsgID,
		ConversationType: int32(msg.MsgData.ConversationType),
		ContentType:      int32(msg.MsgData.ContentType),
		CreateTime:       int64(msg.MsgData.ClientTime),
		Content:          msg.MsgData.Content,
	}
}
func (l *SendMsgLogic) callbackBeforeSendGroupMsg(msg *chatpb.SendMsgReq) (canSend bool, err error) {
	if !l.svcCtx.Config.Callback.CallbackBeforeSendGroupMsg.Enable {
		return true, nil
	}
	commonCallbackReq := l.copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = msgcallbackpb.CallbackCommand_BeforeSendGroupMsg
	req := msgcallbackpb.CallbackSendGroupMsgReq{
		CommonCallbackReq: &commonCallbackReq,
		GroupID:           msg.MsgData.ReceiveID,
	}
	resp, err := l.svcCtx.MsgCallback.CallbackBeforeSendGroupMsg(l.ctx, &req)
	if err != nil {
		if l.svcCtx.Config.Callback.CallbackBeforeSendGroupMsg.ContinueOnError {
			return true, err
		} else {
			return false, err
		}
	}
	if resp.ActionCode == msgcallbackpb.ActionCode_Forbidden && resp.ErrCode == msgcallbackpb.ErrCode_HandleSuccess {
		return false, nil
	}
	return true, err
}

func (l *SendMsgLogic) callbackAfterSendGroupMsg(msg *chatpb.SendMsgReq) error {
	if !l.svcCtx.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return nil
	}
	commonCallbackReq := l.copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = msgcallbackpb.CallbackCommand_AfterSendGroupMsg
	req := msgcallbackpb.CallbackSendGroupMsgReq{
		CommonCallbackReq: &commonCallbackReq,
		GroupID:           msg.MsgData.ReceiveID,
	}
	_, err := l.svcCtx.MsgCallback.CallbackAfterSendGroupMsg(l.ctx, &req)
	if err != nil {
		return err
	}
	return nil
}
