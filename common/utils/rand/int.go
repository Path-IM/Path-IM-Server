package randUtils

import (
	"math/rand"
	"time"
)

// RandInt 随机int
func RandInt(i int, j int) int {
	rand.Seed(time.Now().UnixNano())
	return i + rand.Intn(j-i)
}

func RandInt64(i int64, j int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return i + rand.Int63n(j-i)
}
