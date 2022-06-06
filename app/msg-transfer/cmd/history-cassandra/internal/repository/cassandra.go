package repository

import (
	"context"
	"fmt"
	chatpb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"time"
)

const singleGocMsgNum = 5000
const groupGocMsgNum = 5000

func indexGen(uid string, seqSuffix uint32) string {
	return uid + ":" + strconv.FormatInt(int64(seqSuffix), 10)
}

func getSeqUid(uid string, seq uint32) string {
	seqSuffix := seq / singleGocMsgNum
	return indexGen(uid, seqSuffix)
}
func getSeqGroupId(groupId string, seq uint32) string {
	seqSuffix := seq / groupGocMsgNum
	return indexGen(groupId, seqSuffix)
}

func (r *Rep) SaveUserChatCassandra2(spanCtx context.Context, uid string, sendTime int64, m *chatpb.MsgDataToDB) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.svcCtx.Config.Cassandra.TimeoutSecond)*time.Second)
	seqUid := getSeqUid(uid, m.MsgData.Seq)
	{
		err := r.CassandraClient.Query(
			fmt.Sprintf(
				`UPDATE %s.%s SET msgs = msgs + ? WHERE uid = ?`,
				r.svcCtx.Config.Cassandra.Keyspace,
				r.svcCtx.Config.Cassandra.SingleChatMsgTableName,
			),
			[]map[int64][]byte{{
				sendTime: m.MsgData.Bytes(),
			}}, seqUid,
		).WithContext(ctx).Exec()
		if err != nil {
			logx.WithContext(spanCtx).Errorf("SaveUserChatCassandra2 error: %v", err)
			return err
		}
	}
	return nil
}

func (r *Rep) SaveGroupChatCassandra(spanCtx context.Context, groupId string, sendTime int64, m *chatpb.MsgDataToDB) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.svcCtx.Config.Cassandra.TimeoutSecond)*time.Second)
	seqGroupId := getSeqGroupId(groupId, m.MsgData.Seq)
	{
		err := r.CassandraClient.Query(
			fmt.Sprintf(
				`UPDATE %s.%s SET msgs = msgs + ? WHERE groupid = ?`,
				r.svcCtx.Config.Cassandra.Keyspace,
				r.svcCtx.Config.Cassandra.GroupChatMsgTableName,
			),
			[]map[int64][]byte{{
				sendTime: m.MsgData.Bytes(),
			}}, seqGroupId,
		).WithContext(ctx).Exec()
		if err != nil {
			logx.WithContext(spanCtx).Errorf("SaveGroupChatCassandra error: %v", err)
			return err
		}
	}
	return nil
}
