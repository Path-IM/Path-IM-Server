package push

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/showurl/Zero-IM-Server/app/msg-push/cmd/rpc/internal/config"
	"github.com/showurl/Zero-IM-Server/app/msg-push/cmd/rpc/internal/sdk/jpush/common"
	"github.com/showurl/Zero-IM-Server/app/msg-push/cmd/rpc/internal/sdk/jpush/requestBody"
	"io/ioutil"
	"net/http"
)

type JPush struct {
	Config config.Config
}

func (j *JPush) Push(ctx context.Context, accounts []string, alert, detailContent string) (string, error) {
	if accounts == nil || len(accounts) == 0 {
		return "", nil
	}
	var pf requestBody.Platform
	pf.SetAll()
	var au requestBody.Audience
	au.SetAlias(accounts)
	var no requestBody.Notification
	no.SetAlert(alert, j.Config.Jpns.PushIntent)
	var me requestBody.Message
	me.SetMsgContent(detailContent)
	var o requestBody.Options
	o.SetApnsProduction(false)
	var po requestBody.PushObj
	po.SetPlatform(&pf)
	po.SetAudience(&au)
	po.SetNotification(&no)
	po.SetMessage(&me)
	po.SetOptions(&o)

	con, err := json.Marshal(po)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", j.Config.Jpns.PushUrl, bytes.NewBuffer(con))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", common.GetAuthorization(j.Config.Jpns.AppKey, j.Config.Jpns.MasterSecret))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
