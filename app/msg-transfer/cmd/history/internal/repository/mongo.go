package repository

import (
	"context"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/showurl/Zero-IM-Server/app/msg-transfer/cmd/history/model"
	chatpb "github.com/showurl/Zero-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"time"
)

const singleGocMsgNum = 5000
const superGroupGocMsgNum = 5000

func indexGen(uid string, seqSuffix uint32) string {
	return uid + ":" + strconv.FormatInt(int64(seqSuffix), 10)
}

func getSeqUid(uid string, seq uint32) string {
	seqSuffix := seq / singleGocMsgNum
	return indexGen(uid, seqSuffix)
}
func getSeqGroupId(groupId string, seq uint32) string {
	seqSuffix := seq / superGroupGocMsgNum
	return indexGen(groupId, seqSuffix)
}

func (r *Rep) SaveUserChatMongo2(spanCtx context.Context, uid string, sendTime int64, m *chatpb.MsgDataToDB) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.svcCtx.Config.Mongo.DBTimeout)*time.Second)
	c := r.MongoClient.Database(r.svcCtx.Config.Mongo.DBDatabase).Collection(r.svcCtx.Config.Mongo.SingleChatMsgCollectionName)
	seqUid := getSeqUid(uid, m.MsgData.Seq)
	filter := bson.M{"uid": seqUid}
	var err error
	sMsg := model.MsgInfo{}
	sMsg.SendTime = sendTime
	if sMsg.Msg, err = proto.Marshal(m.MsgData); err != nil {
		logx.WithContext(spanCtx).Errorf("proto.Marshal error: %s", err.Error())
		return err
	}
	err = c.FindOneAndUpdate(ctx, filter, bson.M{"$push": bson.M{"msg": sMsg}}).Err()
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			logx.WithContext(spanCtx).Errorf("SaveUserChatMongo2 error: %s", err.Error())
			return err
		}
		sChat := model.UserChat{}
		sChat.UID = seqUid
		sChat.Msg = append(sChat.Msg, sMsg)
		if _, err = c.InsertOne(ctx, &sChat); err != nil {
			logx.WithContext(spanCtx).Errorf("InsertOne error: %s", err.Error())
			return err
		}
	}
	return nil
}

func (r *Rep) SaveSuperGroupChatMongo2(spanCtx context.Context, groupId string, sendTime int64, m *chatpb.MsgDataToDB) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.svcCtx.Config.Mongo.DBTimeout)*time.Second)
	c := r.MongoClient.Database(r.svcCtx.Config.Mongo.DBDatabase).Collection(r.svcCtx.Config.Mongo.SuperGroupChatMsgCollectionName)
	seqGroupId := getSeqGroupId(groupId, m.MsgData.Seq)
	filter := bson.M{"groupid": seqGroupId}
	var err error
	sMsg := model.MsgInfo{}
	sMsg.SendTime = sendTime
	if sMsg.Msg, err = proto.Marshal(m.MsgData); err != nil {
		logx.WithContext(spanCtx).Errorf("proto.Marshal error: %s", err.Error())
		return err
	}
	err = c.FindOneAndUpdate(ctx, filter, bson.M{"$push": bson.M{"msg": sMsg}}).Err()
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			logx.WithContext(spanCtx).Errorf("SaveSuperGroupChatMongo2 error: %s", err.Error())
			return err
		}
		sChat := model.GroupChat{}
		sChat.GroupID = seqGroupId
		sChat.Msg = append(sChat.Msg, sMsg)
		if _, err = c.InsertOne(ctx, &sChat); err != nil {
			logx.WithContext(spanCtx).Errorf("InsertOne error: %s", err.Error())
			return err
		}
	}
	return nil
}
