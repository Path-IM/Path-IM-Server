package wsrepository

import "testing"

func TestRedisHgetall(t *testing.T) {
	redis := newRedis("127.0.0.1:6379", "123456", "node", false)
	hgetall, err := redis.Hgetall("123")
	if err != nil {
		t.Error(err)
	}
	t.Log(hgetall)
}
