// Code generated by goctl. DO NOT EDIT!
// Source: chat.proto

package chat

import (
	"context"

	"github.com/showurl/Path-IM-Server/app/msg/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	GetMaxAndMinSeqReq                = pb.GetMaxAndMinSeqReq
	GetMaxAndMinSeqResp               = pb.GetMaxAndMinSeqResp
	GetMaxAndMinSuperGroupSeqReq      = pb.GetMaxAndMinSuperGroupSeqReq
	GetMaxAndMinSuperGroupSeqResp     = pb.GetMaxAndMinSuperGroupSeqResp
	GetMaxAndMinSuperGroupSeqRespItem = pb.GetMaxAndMinSuperGroupSeqRespItem
	MsgDataToDB                       = pb.MsgDataToDB
	MsgDataToMQ                       = pb.MsgDataToMQ
	PullMessageBySuperGroupSeqListReq = pb.PullMessageBySuperGroupSeqListReq
	PushMsgDataToMQ                   = pb.PushMsgDataToMQ
	PushMsgToSuperGroupDataToMQ       = pb.PushMsgToSuperGroupDataToMQ
	SendMsgReq                        = pb.SendMsgReq
	SendMsgResp                       = pb.SendMsgResp
	WrapPullMessageBySeqListReq       = pb.WrapPullMessageBySeqListReq
	WrapPullMessageBySeqListResp      = pb.WrapPullMessageBySeqListResp

	Chat interface {
		GetMaxAndMinSeq(ctx context.Context, in *GetMaxAndMinSeqReq, opts ...grpc.CallOption) (*GetMaxAndMinSeqResp, error)
		GetSuperGroupMaxAndMinSeq(ctx context.Context, in *GetMaxAndMinSuperGroupSeqReq, opts ...grpc.CallOption) (*GetMaxAndMinSuperGroupSeqResp, error)
		PullMessageBySeqList(ctx context.Context, in *WrapPullMessageBySeqListReq, opts ...grpc.CallOption) (*WrapPullMessageBySeqListResp, error)
		PullMessageBySuperGroupSeqList(ctx context.Context, in *PullMessageBySuperGroupSeqListReq, opts ...grpc.CallOption) (*WrapPullMessageBySeqListResp, error)
		SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgResp, error)
	}

	defaultChat struct {
		cli zrpc.Client
	}
)

func NewChat(cli zrpc.Client) Chat {
	return &defaultChat{
		cli: cli,
	}
}

func (m *defaultChat) GetMaxAndMinSeq(ctx context.Context, in *GetMaxAndMinSeqReq, opts ...grpc.CallOption) (*GetMaxAndMinSeqResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.GetMaxAndMinSeq(ctx, in, opts...)
}

func (m *defaultChat) GetSuperGroupMaxAndMinSeq(ctx context.Context, in *GetMaxAndMinSuperGroupSeqReq, opts ...grpc.CallOption) (*GetMaxAndMinSuperGroupSeqResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.GetSuperGroupMaxAndMinSeq(ctx, in, opts...)
}

func (m *defaultChat) PullMessageBySeqList(ctx context.Context, in *WrapPullMessageBySeqListReq, opts ...grpc.CallOption) (*WrapPullMessageBySeqListResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.PullMessageBySeqList(ctx, in, opts...)
}

func (m *defaultChat) PullMessageBySuperGroupSeqList(ctx context.Context, in *PullMessageBySuperGroupSeqListReq, opts ...grpc.CallOption) (*WrapPullMessageBySeqListResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.PullMessageBySuperGroupSeqList(ctx, in, opts...)
}

func (m *defaultChat) SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.SendMsg(ctx, in, opts...)
}
