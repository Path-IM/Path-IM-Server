package xmgo

import (
	"context"
	"github.com/showurl/Zero-IM-Server/common/xmgo/global"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetClient(
	cfg global.MongoConfig,
) *mongo.Client {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.Uri))
	if err != nil {
		panic(err)
	}
	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	return mongoClient
}
