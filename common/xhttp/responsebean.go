package xhttp

type (
	NullJson struct{}

	ResponseSuccessBean struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
)

func Success(data interface{}) *ResponseSuccessBean {
	return &ResponseSuccessBean{200, "OK", data}
}

type ResponseErrorBean struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Error(errCode int, errMsg string) *ResponseErrorBean {
	return &ResponseErrorBean{errCode, errMsg}
}

func NewResp(resp interface{}) *ResponseSuccessBean {
	return &ResponseSuccessBean{
		Data: resp,
	}
}
