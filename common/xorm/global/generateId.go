package global

import (
	"github.com/bwmarrin/snowflake"
	"github.com/showurl/Zero-IM-Server/common/utils/encrypt"
	"math/rand"
	"time"
)

var (
	c           = make(chan string)
	nd          *snowflake.Node
	NodId       int64
	uniquePodId string
)

func init() {
	rand.Seed(time.Now().UnixNano())
	NodId = int64(rand.Intn(128))
	nd, _ = snowflake.NewNode(NodId)
	uniquePodId = nd.Generate().String()
	go loop()
}

func loop() {
	for {
		c <- nd.Generate().String()
	}
}

// GetID 获取ID
// 参数: podId 当前机器标识 这里使用 pod ip
func GetID() string {
	md5 := encrypt.Md5(uniquePodId + <-c)
	//fmt.Println("md5:", md5)
	return md5
}

func ReplacePodID(podId string) {
	uniquePodId = podId
}
