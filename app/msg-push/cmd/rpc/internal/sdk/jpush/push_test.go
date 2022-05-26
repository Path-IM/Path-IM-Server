package push

import (
	"context"
	"github.com/showurl/Path-IM-Server/app/msg-push/cmd/rpc/internal/config"
	"testing"
)

func TestPush(t *testing.T) {
	push := &JPush{Config: config.Config{
		Jpns: config.JpnsConf{
			PushIntent:   `intent:#Intent;component=com.xxx.xxx/com.xxx.xxx.MainActivity;end`,
			PushUrl:      `https://api.jpush.cn/v3/push`,
			AppKey:       ``,
			MasterSecret: ``,
		},
	}}
	resp, err := push.Push(context.Background(), []string{
		"xxx",
	}, "代码-你根本没在沈阳", "代码-我是云南的")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
