package global

type MysqlConfig struct {
	Addr         string
	MaxIdleConns int
	MaxOpenConns int
	LogLevel     string
}
