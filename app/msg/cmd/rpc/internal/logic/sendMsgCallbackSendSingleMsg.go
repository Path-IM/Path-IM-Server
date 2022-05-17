package logic

import (
	msgcallbackpb "github.com/showurl/Zero-IM-Server/app/msg-callback/cmd/rpc/pb"
	chatpb "github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"
)

func (l *SendMsgLogic) callbackBeforeSendSingleMsg(msg *chatpb.SendMsgReq) (canSend bool, err error) {
	if !l.svcCtx.Config.Callback.CallbackBeforeSendSingleMsg.Enable {
		return true, nil
	}

	commonCallbackReq := l.copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = msgcallbackpb.CallbackCommand_BeforeSendSingleMsg
	req := msgcallbackpb.CallbackSendSingleMsgReq{
		CommonCallbackReq: &commonCallbackReq,
		RecvID:            msg.MsgData.RecvID,
	}
	resp, err := l.svcCtx.MsgCallback.CallbackBeforeSendSingleMsg(l.ctx, &req)
	if err != nil {
		if l.svcCtx.Config.Callback.CallbackBeforeSendSingleMsg.ContinueOnError {
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

func (l *SendMsgLogic) callbackAfterSendSingleMsg(msg *chatpb.SendMsgReq) error {

	if !l.svcCtx.Config.Callback.CallbackAfterSendSingleMsg.Enable {
		return nil
	}
	commonCallbackReq := l.copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = msgcallbackpb.CallbackCommand_AfterSendSingleMsg
	req := msgcallbackpb.CallbackSendSingleMsgReq{
		CommonCallbackReq: &commonCallbackReq,
		RecvID:            msg.MsgData.RecvID,
	}
	_, err := l.svcCtx.MsgCallback.CallbackAfterSendSingleMsg(l.ctx, &req)
	if err != nil {
		return err
	}
	return nil
}
