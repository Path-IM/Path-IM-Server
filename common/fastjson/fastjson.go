package fastjson

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
)

var fJson = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return fJson.Unmarshal(data, v)
}
