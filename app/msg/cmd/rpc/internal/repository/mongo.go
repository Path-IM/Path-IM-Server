package repository

import (
	"context"
	"github.com/Path-IM/Path-IM-Server/app/msg-transfer/model"
	"github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	numUtils "github.com/Path-IM/Path-IM-Server/common/utils/num"
	"github.com/golang/protobuf/proto"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *MongoHistory) GetMsgBySeqList(uid string, seqList []uint32) (seqMsg []*pb.MsgData, err error) {
	var hasSeqList []uint32
	singleCount := 0
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.svcCtx.Config.Mongo.DBTimeout)*time.Second)
	c := r.MongoClient.Database(r.svcCtx.Config.Mongo.DBDatabase).Collection(r.svcCtx.Config.Mongo.SingleChatMsgCollectionName)
	m := func(uid string, seqList []uint32) map[string][]uint32 {
		t := make(map[string][]uint32)
		for i := 0; i < len(seqList); i++ {
			seqUid := getSeqUid(uid, seqList[i])
			if value, ok := t[seqUid]; !ok {
				var temp []uint32
				t[seqUid] = append(temp, seqList[i])
			} else {
				t[seqUid] = append(value, seqList[i])
			}
		}
		return t
	}(uid, seqList)
	sChat := model.UserChat{}
	for seqUid, value := range m {
		if err = c.FindOne(ctx, bson.M{"uid": seqUid}).Decode(&sChat); err != nil {
			logx.Error("findone uid failed:", err.Error())
			continue
		}
		singleCount = 0
		for i := 0; i < len(sChat.Msg); i++ {
			msg := new(pb.MsgData)
			if err = proto.Unmarshal(sChat.Msg[i].Msg, msg); err != nil {
				logx.Error("Unmarshal failed:", err.Error())
				return nil, err
			}
			if numUtils.IsContainUInt32(msg.Seq, value) {
				seqMsg = append(seqMsg, msg)
				hasSeqList = append(hasSeqList, msg.Seq)
				singleCount++
				if singleCount == len(value) {
					break
				}
			}
		}
	}
	if len(hasSeqList) != len(seqList) {
		var diff []uint32
		diff = numUtils.DifferenceUInt32(hasSeqList, seqList)
		exceptionMSg := genExceptionMessageBySeqList(diff)
		seqMsg = append(seqMsg, exceptionMSg...)

	}
	return seqMsg, nil
}

func (r *MongoHistory) GetMsgBySuperGroupSeqList(groupId string, seqList []uint32) (seqMsg []*pb.MsgData, err error) {
	var hasSeqList []uint32
	singleCount := 0
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.svcCtx.Config.Mongo.DBTimeout)*time.Second)
	c := r.MongoClient.Database(r.svcCtx.Config.Mongo.DBDatabase).Collection(r.svcCtx.Config.Mongo.SuperGroupChatMsgCollectionName)
	m := func(groupId string, seqList []uint32) map[string][]uint32 {
		t := make(map[string][]uint32)
		for i := 0; i < len(seqList); i++ {
			seqGroupId := getSeqGroupId(groupId, seqList[i])
			if value, ok := t[seqGroupId]; !ok {
				var temp []uint32
				t[seqGroupId] = append(temp, seqList[i])
			} else {
				t[seqGroupId] = append(value, seqList[i])
			}
		}
		return t
	}(groupId, seqList)
	sChat := model.GroupChat{}
	for seqGroupId, value := range m {
		if err = c.FindOne(ctx, bson.M{"groupid": seqGroupId}).Decode(&sChat); err != nil {
			logx.Error("findone uid failed:", err.Error())
			continue
		}
		singleCount = 0
		for i := 0; i < len(sChat.Msg); i++ {
			msg := new(pb.MsgData)
			if err = proto.Unmarshal(sChat.Msg[i].Msg, msg); err != nil {
				logx.Error("Unmarshal failed:", err.Error())
				return nil, err
			}
			if numUtils.IsContainUInt32(msg.Seq, value) {
				seqMsg = append(seqMsg, msg)
				hasSeqList = append(hasSeqList, msg.Seq)
				singleCount++
				if singleCount == len(value) {
					break
				}
			}
		}
	}
	if len(hasSeqList) != len(seqList) {
		var diff []uint32
		diff = numUtils.DifferenceUInt32(hasSeqList, seqList)
		exceptionMSg := genExceptionMessageBySeqList(diff)
		seqMsg = append(seqMsg, exceptionMSg...)

	}
	return seqMsg, nil
}
