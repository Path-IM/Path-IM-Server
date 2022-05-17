package logic

import (
	msgcallbackpb "github.com/showurl/Zero-IM-Server/app/msg-callback/cmd/rpc/pb"
	chatpb "github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"
)

func (l *SendMsgLogic) callbackBeforeSendSuperGroupMsg(msg *chatpb.SendMsgReq) (canSend bool, err error) {
	if !l.svcCtx.Config.Callback.CallbackBeforeSendSuperGroupMsg.Enable {
		return true, nil
	}
	commonCallbackReq := l.copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = msgcallbackpb.CallbackCommand_BeforeSendSuperGroupMsg
	req := msgcallbackpb.CallbackSendSuperGroupMsgReq{
		CommonCallbackReq: &commonCallbackReq,
		SuperGroupID:      msg.MsgData.GroupID,
	}
	resp, err := l.svcCtx.MsgCallback.CallbackBeforeSendSuperGroupMsg(l.ctx, &req)
	if err != nil {
		if l.svcCtx.Config.Callback.CallbackBeforeSendSuperGroupMsg.ContinueOnError {
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

func (l *SendMsgLogic) callbackAfterSendSuperGroupMsg(msg *chatpb.SendMsgReq) error {
	if !l.svcCtx.Config.Callback.CallbackAfterSendSuperGroupMsg.Enable {
		return nil
	}
	commonCallbackReq := l.copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = msgcallbackpb.CallbackCommand_AfterSendSuperGroupMsg
	req := msgcallbackpb.CallbackSendSuperGroupMsgReq{
		CommonCallbackReq: &commonCallbackReq,
		SuperGroupID:      msg.MsgData.GroupID,
	}
	_, err := l.svcCtx.MsgCallback.CallbackAfterSendSuperGroupMsg(l.ctx, &req)
	if err != nil {
		return err
	}
	return nil
}
