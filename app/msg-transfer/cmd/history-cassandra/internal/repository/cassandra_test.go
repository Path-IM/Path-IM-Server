package repository

import (
	"fmt"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/Path-IM/Path-IM-Server/common/xcql"
	"github.com/gocql/gocql"
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
	"time"
)

var (
	dt = &pb.MsgData{
		SendID:           "12345678987654321234567",
		RecvID:           "asdfghjkjhgfdsasdfghj",
		GroupID:          "q2w34rtyuytredswasdefgh",
		ClientMsgID:      "2345yutrewertyutrew",
		ServerMsgID:      "hgfdaerwtyuitrewtyu",
		SenderPlatformID: 2,
		SenderNickname:   "213456ytrewertyutrew",
		SenderFaceURL:    "213456yutrewrtyutrewrtyrewrtytrewrtyutrew",
		SessionType:      51,
		MsgFrom:          12,
		ContentType:      1615,
		Content: []byte(`{
			"msg_type": "text",
			"text": "hello world"
}`),
		Seq:        1212,
		SendTime:   1234532345676,
		CreateTime: 23456754323456,
		OfflinePushInfo: &pb.OfflinePushInfo{
			Title:         "123456r4ewqrt",
			Desc:          "2345trewwrt",
			Ex:            "wergfwdqwefgds",
			IOSPushSound:  "ewrgfdwq",
			IOSBadgeCount: true,
		},
		AtUserIDList: []string{"12345678987654321234567", "asdfghjkjhgfdsasdfghj", "q2w34rtyuytredswasdefgh"},
		Options: map[string]bool{
			"option1": true,
			"option2": false,
		},
	}
	bytes = dt.Bytes()
)

func _getSession() *gocql.Session {
	client := xcql.GetClient(xcql.CassandraConfig{
		Hosts:         []string{"172.30.1.1"},
		Port:          9042,
		Keyspace:      "",
		Username:      "cassandra",
		Password:      "cassandra",
		TimeoutSecond: 10,
		Consistency:   "ONE",
	})
	err := client.Query(`DESCRIBE path_im_server`).Exec()
	if err != nil {
		if e, ok := err.(gocql.RequestError); ok {
			if e.Code() == 8704 {
				err = client.Query(`CREATE KEYSPACE path_im_server WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 }`).Exec()
				if err != nil {
					panic(err)
				}
			}
		} else {
			panic(err)
		}
	}
	if err := client.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS path_im_server.%s (seq bigint, send_time bigint, msg_data blob, PRIMARY KEY (seq))", "history")).Exec(); err != nil {
		logx.Errorf("create user chat table error: %s", err)
		panic(err)
	}
	return client
}

func _getSession2() *gocql.Session {
	client := xcql.GetClient(xcql.CassandraConfig{
		Hosts:         []string{"172.30.1.1"},
		Port:          9042,
		Keyspace:      "",
		Username:      "cassandra",
		Password:      "cassandra",
		TimeoutSecond: 10,
		Consistency:   "QUORUM",
	})
	err := client.Query(`DESCRIBE path_im_server`).Exec()
	if err != nil {
		if e, ok := err.(gocql.RequestError); ok {
			if e.Code() == 8704 {
				err = client.Query(`CREATE KEYSPACE path_im_server WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 }`).Exec()
				if err != nil {
					panic(err)
				}
			}
		} else {
			panic(err)
		}
	}
	if err := client.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS path_im_server.%s (seq bigint, msgs list< frozen< map<bigint, blob> > >, PRIMARY KEY (seq))", "history_push")).Exec(); err != nil {
		logx.Errorf("create user chat table error: %s", err)
		panic(err)
	}
	return client
}

/*
goos: darwin
goarch: arm64
pkg: github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-cassandra/internal/repository
BenchmarkInsert
BenchmarkInsert-10    	     447	   2756859 ns/op
PASS
*/
func BenchmarkInsert(b *testing.B) {
	session := _getSession()
	defer session.Close()
	for i := 0; i < b.N; i++ {
		err := session.Query("INSERT INTO path_im_server.history (seq, send_time, msg_data) VALUES (?, ?, ?)",
			i, time.Now().Unix(), bytes,
		).Exec()
		if err != nil {
			logx.Errorf("insert error: %s", err)
			panic(err)
		}
	}
}

/*
goos: darwin
goarch: arm64
pkg: github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-cassandra/internal/repository
BenchmarkPush
BenchmarkPush-10    	     644	   1952197 ns/op
PASS
*/
func BenchmarkPush(b *testing.B) {
	session := _getSession2()
	defer session.Close()
	e := session.Query(`DELETE FROM path_im_server.history_push WHERE seq = ?`, 0).Exec()
	if e != nil {
		logx.Errorf("delete error: %s", e)
		panic(e)
	}
	for i := 0; i < b.N; i++ {
		{
			err := session.Query("UPDATE path_im_server.history_push SET msgs = msgs + ? WHERE seq = ?",
				[]map[int64][]byte{map[int64][]byte{
					time.Now().UnixMilli(): bytes,
				}}, 3,
			).Exec()
			if err != nil {
				logx.Errorf("insert error: %s", err)
				panic(err)
			}
		}
	}
}

func BenchmarkList(b *testing.B) {

}

type model struct {
	Seq  int64
	Msgs []map[int64][]byte
}

/*
goos: darwin
goarch: arm64
pkg: github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-cassandra/internal/repository
BenchmarkFirst
BenchmarkFirst-10    	     380	   2774049 ns/op
PASS
// update
goos: darwin
goarch: arm64
pkg: github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-cassandra/internal/repository
BenchmarkFirst
BenchmarkFirst-10    	      12	  95737656 ns/op
PASS
*/
func BenchmarkFirst(b *testing.B) {
	session := _getSession()
	defer session.Close()
	for i := 0; i < b.N; i++ {
		m := &model{}
		err := session.Query("SELECT seq, msgs FROM path_im_server.history_push WHERE seq = ?",
			2,
		).Scan(&m.Seq, &m.Msgs)
		if err != nil {
			logx.Errorf("select error: %s", err)
			panic(err)
		}
	}
}
