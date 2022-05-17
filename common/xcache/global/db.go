package global

type DB int

func (d DB) Int() int {
	return int(d)
}

const (
	DB0Default DB = iota // 默认
	DB1Token             // Token
	DB2String            // String
	DB2ZSet              // ZSet
)
