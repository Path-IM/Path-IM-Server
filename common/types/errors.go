package types

import "github.com/Path-IM/Path-IM-Server/common/xerr"

var WSDataError = xerr.New(3001, "ws data error")

// error code
const (
	ErrCodeOK     = iota // 成功
	ErrCodeFailed        // 失败
	ErrCodeLimit         // 限流

	ErrCodeProtoUnmarshal = 400 // proto解析错误
	ErrCodeParams         = 401 // 参数错误
)
