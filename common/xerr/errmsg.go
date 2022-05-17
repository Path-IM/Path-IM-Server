package xerr

var message map[int]string

func init() {
	message = make(map[int]string)
	message[OK] = "SUCCESS"
	message[BAD_REUQEST_ERROR] = "服务器繁忙,请稍后再试"
	message[REUQES_PARAM_ERROR] = "参数错误"
}

func MapErrMsg(errcode int) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "服务器繁忙,请稍后再试"
	}
}
