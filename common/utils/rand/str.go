package randUtils

import (
	"math/rand"
	"strconv"
	"time"
)

// RandPhone 随机手机号
func RandPhone() string {
	return "13" + RandNum(8)
}

// RandNum 随机数字
func RandNum(i int) string {
	var str string
	for j := 0; j < i; j++ {
		str += strconv.Itoa(RandInt(0, 9))
	}
	return str
}

// RandString 随机字符串
func RandString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
