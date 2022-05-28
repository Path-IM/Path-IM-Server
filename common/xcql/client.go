package xcql

import (
	"github.com/gocql/gocql"
	"time"
)

func GetClient(conf CassandraConfig) *gocql.Session {
	cluster := gocql.NewCluster(conf.Hosts...)
	//cluster.Keyspace = conf.Keyspace
	//cluster.Consistency = gocql.Quorum
	cluster.Consistency = gocql.ParseConsistency(conf.Consistency)
	cluster.Timeout = time.Second * time.Duration(conf.TimeoutSecond)
	cluster.ProtoVersion = 4
	cluster.Port = conf.Port
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: conf.Username,
		Password: conf.Password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	return session
}
