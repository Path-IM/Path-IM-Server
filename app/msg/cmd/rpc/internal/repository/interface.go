package repository

import (
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/internal/svc"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/mongo"
)

type IPullHistoryMsg interface {
	GetMsgBySeqList(uid string, seqList []uint32) (seqMsg []*pb.MsgData, err error)
	GetMsgBySuperGroupSeqList(groupId string, seqList []uint32) (seqMsg []*pb.MsgData, err error)
}
type MongoHistory struct {
	svcCtx      *svc.ServiceContext
	MongoClient *mongo.Client
}

type CassandraHistory struct {
	svcCtx          *svc.ServiceContext
	CassandraClient *gocql.Session
}
