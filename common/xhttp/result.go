package xhttp

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/showurl/Zero-IM-Server/common/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
	"net/http"
)

// HttpResult http方法
func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	go func() {
		logx.WithContext(r.Context()).Infof(
			"[X-Real-IP][%s] [User-Agent][%s] ",
			r.Header.Get("x-real-ip"),
			r.Header.Get("user-agent"),
		)
	}()
	if err == nil {
		// 成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		// 错误返回
		errcode := xerr.BAD_REUQEST_ERROR
		errmsg := "服务器繁忙，请稍后再试"
		if e, ok := err.(*xerr.CodeError); ok {
			// 自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			originErr := errors.Cause(err) // err类型
			if gstatus, ok := status.FromError(originErr); ok {
				// grpc err错误
				errmsg = gstatus.Message()
			}
		}
		logx.WithContext(r.Context()).Errorf("【GATEWAY-SRV-ERR】 : %+v ", err)

		httpx.WriteJson(w, http.StatusOK, Error(errcode, errmsg))
	}
}

// ParamErrorResult http参数错误返回
func ParamErrorResult(r *http.Request, w http.ResponseWriter, err error) {
	errMsg := fmt.Sprintf("%s ,%s", xerr.MapErrMsg(xerr.REUQES_PARAM_ERROR), err.Error())
	httpx.WriteJson(w, http.StatusBadRequest, Error(xerr.REUQES_PARAM_ERROR, errMsg))
}
