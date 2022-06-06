package types

// error code
const (
	ErrCodeOK     = iota // 成功
	ErrCodeFailed        // 失败
	ErrCodeLimit         // 限流

	ErrCodeProtoUnmarshal = 400 // proto解析错误
	ErrCodeParams         = 401 // 参数错误
)
