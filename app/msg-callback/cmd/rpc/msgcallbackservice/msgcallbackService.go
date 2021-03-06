// Code generated by goctl. DO NOT EDIT!
// Source: msgcallback.proto

package msgcallbackservice

import (
	"context"

	"github.com/Path-IM/Path-IM-Server/app/msg-callback/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CallbackSendGroupMsgReq  = pb.CallbackSendGroupMsgReq
	CallbackSendSingleMsgReq = pb.CallbackSendSingleMsgReq
	CommonCallbackReq        = pb.CommonCallbackReq
	CommonCallbackResp       = pb.CommonCallbackResp

	MsgcallbackService interface {
		CallbackBeforeSendGroupMsg(ctx context.Context, in *CallbackSendGroupMsgReq, opts ...grpc.CallOption) (*CommonCallbackResp, error)
		CallbackAfterSendGroupMsg(ctx context.Context, in *CallbackSendGroupMsgReq, opts ...grpc.CallOption) (*CommonCallbackResp, error)
		CallbackBeforeSendSingleMsg(ctx context.Context, in *CallbackSendSingleMsgReq, opts ...grpc.CallOption) (*CommonCallbackResp, error)
		CallbackAfterSendSingleMsg(ctx context.Context, in *CallbackSendSingleMsgReq, opts ...grpc.CallOption) (*CommonCallbackResp, error)
	}

	defaultMsgcallbackService struct {
		cli zrpc.Client
	}
)

func NewMsgcallbackService(cli zrpc.Client) MsgcallbackService {
	return &defaultMsgcallbackService{
		cli: cli,
	}
}

func (m *defaultMsgcallbackService) CallbackBeforeSendGroupMsg(ctx context.Context, in *CallbackSendGroupMsgReq, opts ...grpc.CallOption) (*CommonCallbackResp, error) {
	client := pb.NewMsgcallbackServiceClient(m.cli.Conn())
	return client.CallbackBeforeSendGroupMsg(ctx, in, opts...)
}

func (m *defaultMsgcallbackService) CallbackAfterSendGroupMsg(ctx context.Context, in *CallbackSendGroupMsgReq, opts ...grpc.CallOption) (*CommonCallbackResp, error) {
	client := pb.NewMsgcallbackServiceClient(m.cli.Conn())
	return client.CallbackAfterSendGroupMsg(ctx, in, opts...)
}

func (m *defaultMsgcallbackService) CallbackBeforeSendSingleMsg(ctx context.Context, in *CallbackSendSingleMsgReq, opts ...grpc.CallOption) (*CommonCallbackResp, error) {
	client := pb.NewMsgcallbackServiceClient(m.cli.Conn())
	return client.CallbackBeforeSendSingleMsg(ctx, in, opts...)
}

func (m *defaultMsgcallbackService) CallbackAfterSendSingleMsg(ctx context.Context, in *CallbackSendSingleMsgReq, opts ...grpc.CallOption) (*CommonCallbackResp, error) {
	client := pb.NewMsgcallbackServiceClient(m.cli.Conn())
	return client.CallbackAfterSendSingleMsg(ctx, in, opts...)
}
