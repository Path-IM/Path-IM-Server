package testUtils

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

type Do func(r *resty.Request) (*resty.Response, error)

func ParseResponseJsonBody(f Do, resp interface{}) {
	if f != nil {
		r := resty.New().R().SetHeader("Content-Type", "application/json")
		response, err := f(r)
		if err != nil {
			logx.Info(err)
		} else {
			reqBody := response.Request.Body
			logx.Info("=== Start:reqBody")
			buf, err := json.Marshal(reqBody)
			if err != nil {
				logx.Info(reqBody)
			} else {
				var bb bytes.Buffer
				_ = json.Indent(&bb, buf, "", "    ")
				logx.Info(bb.String())
			}
			logx.Info("=== End:reqBody")
			if response == nil {
				logx.Info("nil")
			} else {
				logx.Info("=== Start:respBody http-status：", response.Status())
				var bb bytes.Buffer
				src := []byte(response.String())
				if resp != nil {
					err2 := json.Unmarshal(src, resp)
					if err2 != nil {
						logx.Info("json反序列化失败")
					}
				}
				err = json.Indent(&bb, src, "", "    ")
				if err == nil {
					body := bb.String()
					logx.Info(body)
				} else {
					logx.Info(response.String() + " : " + r.URL)
				}
				logx.Info("=== End:respBody")
			}
		}
	} else {
		logx.Info("f is nil")
	}
}
