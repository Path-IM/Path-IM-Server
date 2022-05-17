package logic

import (
	msgcallbackpb "github.com/showurl/Zero-IM-Server/app/msg-callback/cmd/rpc/pb"
	chatpb "github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"
)

func (l *SendMsgLogic) callbackBeforeSendGroupMsg(msg *chatpb.SendMsgReq) (canSend bool, err error) {
	if !l.svcCtx.Config.Callback.CallbackBeforeSendGroupMsg.Enable {
		return true, nil
	}
	commonCallbackReq := l.copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = msgcallbackpb.CallbackCommand_BeforeSendGroupMsg
	req := msgcallbackpb.CallbackSendGroupMsgReq{
		CommonCallbackReq: &commonCallbackReq,
		GroupID:           msg.MsgData.GroupID,
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
		GroupID:           msg.MsgData.GroupID,
	}
	_, err := l.svcCtx.MsgCallback.CallbackAfterSendGroupMsg(l.ctx, &req)
	if err != nil {
		return err
	}
	return nil
}
