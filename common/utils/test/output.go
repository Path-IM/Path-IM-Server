package testUtils

import (
	"bytes"
	"encoding/json"
	"github.com/showurl/Zero-IM-Server/common/fastjson"
)

func OutputJson(v interface{}) string {
	buf, _ := fastjson.Marshal(v)
	var bb bytes.Buffer
	_ = json.Indent(&bb, buf, "", "    ")
	s := bb.String()
	//fmt.Println(s)
	return s
}
