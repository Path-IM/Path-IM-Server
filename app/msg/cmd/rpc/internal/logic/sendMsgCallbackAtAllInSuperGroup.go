package logic

import (
	msgcallbackpb "github.com/showurl/Zero-IM-Server/app/msg-callback/cmd/rpc/pb"
	chatpb "github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/showurl/Zero-IM-Server/common/types"
)

func (l *SendMsgLogic) callbackAtAllInSuperGroup(msg *chatpb.SendMsgReq) (canSend bool, err error) {
	if !l.svcCtx.Config.Callback.CallbackAtAllInSuperGroup.Enable || msg.MsgData.ContentType != types.Text {
		return true, nil
	}
	commonCallbackReq := l.copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = msgcallbackpb.CallbackCommand_AtAllInSuperGroup
	req := msgcallbackpb.CallbackAtAllInSuperGroupReq{
		CommonCallbackReq: &commonCallbackReq,
		SuperGroupID:      msg.MsgData.GroupID,
	}
	resp, err := l.svcCtx.MsgCallback.CallbackAtAllInSuperGroup(l.ctx, &req)
	if err != nil {
		if l.svcCtx.Config.Callback.CallbackAtAllInSuperGroup.ContinueOnError {
			return true, err
		} else {
			return false, err
		}
	}
	if resp.CommonCallbackResp.ActionCode == msgcallbackpb.ActionCode_Forbidden && resp.CommonCallbackResp.ErrCode == msgcallbackpb.ErrCode_HandleSuccess {
		return false, nil
	}
	return true, err
}
