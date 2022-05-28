package repository

import (
	"fmt"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/cmd/history-cassandra/internal/svc"
	"github.com/Path-IM/Path-IM-Server/common/xcache"
	"github.com/Path-IM/Path-IM-Server/common/xcache/global"
	"github.com/Path-IM/Path-IM-Server/common/xcql"
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
	"github.com/zeromicro/go-zero/core/logx"
)

type Rep struct {
	svcCtx *svc.ServiceContext
	Cache  redis.UniversalClient
	//mysql       *gorm.DB
	CassandraClient *gocql.Session
}

var rep *Rep

func NewRep(svcCtx *svc.ServiceContext) *Rep {
	if rep != nil {
		return rep
	}
	rep = &Rep{
		svcCtx: svcCtx,
		Cache:  xcache.GetClient(svcCtx.Config.Redis.Conf, global.DB(svcCtx.Config.Redis.DB)),
		//mysql:       xorm.GetClient(svcCtx.Config.Mysql),
		CassandraClient: xcql.GetClient(svcCtx.Config.Cassandra.CassandraConfig),
	}
	// 检查 cassandradb 表 和 索引
	rep.CheckTable()
	return rep
}

func (r *Rep) CheckTable() {
	err := r.CassandraClient.Query(`DESCRIBE ` + r.svcCtx.Config.Cassandra.Keyspace).Exec()
	if err != nil {
		if e, ok := err.(gocql.RequestError); ok {
			if e.Code() == 8704 {
				err = r.CassandraClient.Query(`CREATE KEYSPACE ` + r.svcCtx.Config.Cassandra.Keyspace + ` WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 }`).Exec()
				if err != nil {
					panic(err)
				}
			}
		} else {
			panic(err)
		}
	}
	cql1 := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (uid varchar, msgs list< frozen< map<bigint, blob> > >, PRIMARY KEY (uid))", r.svcCtx.Config.Cassandra.Keyspace, r.svcCtx.Config.Cassandra.SingleChatMsgTableName)
	logx.Infof("cql1: %s", cql1)
	if err := r.CassandraClient.Query(cql1).Exec(); err != nil {
		logx.Errorf("create user chat table error: %s", err)
		panic(err)
	}
	cql2 := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (groupid varchar, msgs list< frozen< map<bigint, blob> > >, PRIMARY KEY (groupid))", r.svcCtx.Config.Cassandra.Keyspace, r.svcCtx.Config.Cassandra.SuperGroupChatMsgTableName)
	logx.Infof("cql2: %s", cql2)
	if err := r.CassandraClient.Query(cql2).Exec(); err != nil {
		logx.Errorf("create user chat table error: %s", err)
		panic(err)
	}
}
