package svc

import (
	"github.com/Path-IM/Path-IM-Server/app/im-user/cmd/rpc/imuserservice"
	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/internal/config"
	"github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/internal/sdk"
	push "github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/internal/sdk/jpush"
	mob_push_sdk "github.com/Path-IM/Path-IM-Server/app/msg-push/cmd/rpc/internal/sdk/mobpush"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	offlinePusher sdk.OfflinePusher
	ImUserService imuserservice.ImUserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:        c,
		ImUserService: imuserservice.NewImUserService(zrpc.MustNewClient(c.ImUserRpc)),
	}
	return s
}
func (s *ServiceContext) GetOfflinePusher() sdk.OfflinePusher {
	if s.offlinePusher != nil {
		return s.offlinePusher
	}
	if s.Config.PushType == "jpns" {
		s.offlinePusher = &push.JPush{s.Config}
	} else if s.Config.PushType == "mobpush" {
		{
			if s.Config.MobPush.AppKey == "" {
				panic("mobpush appkey is empty")
			}
			if s.Config.MobPush.AppSecret == "" {
				panic("mobpush appsecret is empty")
			}
		}
		s.offlinePusher = &mob_push_sdk.Pusher{
			AppKey:         s.Config.MobPush.AppKey,
			AppSecret:      s.Config.MobPush.AppSecret,
			ApnsProduction: s.Config.MobPush.ApnsProduction,
			ApnsCateGory:   s.Config.MobPush.ApnsCateGory,
			ApnsSound:      s.Config.MobPush.ApnsSound,
			AndroidSound:   s.Config.MobPush.AndroidSound,
		}
	} else {
		panic("unsupported push type: " + s.Config.PushType)
	}
	return s.offlinePusher
}
