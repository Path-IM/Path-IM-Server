package rc

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Member struct {
	Data interface{}
}

func (myT *Member) Scan(value interface{}) error {
	myT.Data = value
	return nil
}

type Score struct {
	Data interface{}
}

func (myT *Score) Scan(value interface{}) error {
	myT.Data = value
	return nil
}

func (myT Score) Float64() float64 {
	if myT.Data == nil {
		return 0
	}
	switch v := myT.Data.(type) {
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case string:
		parseUint, _ := strconv.ParseUint(hex.EncodeToString([]byte(v)), 16, 32)
		return float64(parseUint)
	case []uint8:
		parseUint, _ := strconv.ParseUint(string(v), 16, 32)
		return float64(parseUint)
	case time.Time:
		return float64(v.UnixNano() / 1e6)
	case *time.Time:
		if v != nil {
			return float64(v.UnixNano() / 1e6)
		} else {
			return 0
		}
	}
	panic(fmt.Sprintf("数据库中的排序字段[%s]无法转换为float64", reflect.TypeOf(myT.Data).String()))
}
