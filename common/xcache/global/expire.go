package global

import (
	"math/rand"
	"time"
)

func ExpireDuration(second int) time.Duration {
	return time.Second * time.Duration(expire(second))
}

func expire(second int) int {
	return second + rand.Intn(second/10)
}

func Expire5Min() int {
	expireTime := 5 * 60
	return expire(expireTime)
}

func Expire10Min() int {
	expireTime := 10 * 60
	return expire(expireTime)
}

func Expire30Min() int {
	expireTime := 30 * 60
	return expire(expireTime)
}

func ExpireHour() int {
	expireTime := 60 * 60
	return expire(expireTime)
}

func ExpireDay() int {
	expireTime := 24 * 60 * 60
	return expire(expireTime)
}
