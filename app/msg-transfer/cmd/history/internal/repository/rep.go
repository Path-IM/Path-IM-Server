package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/showurl/Zero-IM-Server/app/msg-transfer/cmd/history/internal/svc"
	"github.com/showurl/Zero-IM-Server/common/xcache"
	"github.com/showurl/Zero-IM-Server/common/xcache/global"
	"github.com/showurl/Zero-IM-Server/common/xmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Rep struct {
	svcCtx *svc.ServiceContext
	Cache  redis.UniversalClient
	//mysql       *gorm.DB
	MongoClient *mongo.Client
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
		MongoClient: xmgo.GetClient(svcCtx.Config.Mongo.MongoConfig),
	}
	// 检查 mongodb 索引
	rep.CheckMongoIndexSingle()
	rep.CheckMongoIndexSuperGroup()
	return rep
}

func (r *Rep) CheckMongoIndexSingle() {
	type indexResult struct {
		V    int               `bson:"v"`
		Key  map[string]string `bson:"key"`
		Name string            `bson:"name"`
	}
	var results []indexResult
	collection := r.MongoClient.Database(r.svcCtx.Config.Mongo.DBDatabase).Collection(r.svcCtx.Config.Mongo.SingleChatMsgCollectionName)
	cursor, err := collection.Indexes().List(context.Background())
	if err != nil {
		panic(err)
	}
	_ = cursor.All(context.Background(), &results)
	var hasUid = false
	for _, result := range results {
		if len(result.Key) > 0 {
			if indexType, ok := result.Key["uid"]; ok && indexType != "" {
				hasUid = true
				break
			}
		}
	}
	if !hasUid {
		_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys: bson.M{
				"uid": "hashed",
			},
		})
		if err != nil {
			panic(err)
		}
	}
}

func (r *Rep) CheckMongoIndexSuperGroup() {
	type indexResult struct {
		V    int               `bson:"v"`
		Key  map[string]string `bson:"key"`
		Name string            `bson:"name"`
	}
	var results []indexResult
	collection := r.MongoClient.Database(r.svcCtx.Config.Mongo.DBDatabase).Collection(r.svcCtx.Config.Mongo.SuperGroupChatMsgCollectionName)
	cursor, err := collection.Indexes().List(context.Background())
	if err != nil {
		panic(err)
	}
	_ = cursor.All(context.Background(), &results)
	var hasUid = false
	for _, result := range results {
		if len(result.Key) > 0 {
			if indexType, ok := result.Key["groupid"]; ok && indexType != "" {
				hasUid = true
				break
			}
		}
	}
	if !hasUid {
		_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys: bson.M{
				"groupid": "hashed",
			},
		})
		if err != nil {
			panic(err)
		}
	}
}
