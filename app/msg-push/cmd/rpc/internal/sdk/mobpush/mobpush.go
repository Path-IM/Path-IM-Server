package mob_push_sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type Pusher struct {
	AppKey    string
	AppSecret string
	// apns
	ApnsProduction bool
	ApnsCateGory   string
	ApnsSound      string
	// android
	AndroidSound string

	client *PushClient
}

func (p *Pusher) Push(ctx context.Context, userIDList []string, alert, detailContent string) (resp string, err error) {
	if len(userIDList) == 0 {
		return "", nil
	}
	if p.client == nil {
		p.client = NewPushClient(p.AppKey, p.AppSecret)
	}
	model := NewPushModel(p.AppKey)
	// PushNotify
	{
		iosProduct := 1
		if !p.ApnsProduction {
			iosProduct = 0
		}
		var pushMaps []PushMap
		pushMap := make(map[string]interface{})
		e := json.Unmarshal([]byte(detailContent), &pushMap)
		if e == nil {
			for k, v := range pushMap {
				pushMaps = append(pushMaps, PushMap{k, fmt.Sprintf("%v", v)})
			}
		}
		model.PushNotify = PushNotify{
			TaskCron:       0,           // 是否是定时任务：0否，1是
			TaskTime:       0,           // 定时消息 发送时间， taskCron=1时必填 传入时间戳单位（毫秒 例如 1594277916000 ）
			Plats:          []int{1, 2}, // 1 android;2 ios
			IosProduct:     iosProduct,  // plat = 2下，0测试环境，1生产环境，默认1
			OfflineSeconds: 0,           // 离线消息保存时间，如若不传此参数，默认86400s（1天）
			Content:        alert,       // 推送内容, 如果内容长度超过厂商的限制, 则内容会被截断. vivo不允许纯表情
			Title:          "",          // 如果不设置，则默认的通知标题为应用的名称。如果标题长度超过厂商的限制, 则标题会被截断. vivo不允许纯表情
			Type:           1,           // 推送类型：1通知；2自定义
			Url:            "",          // 1 link跳转 moblink功能的的uri
			CustomNotify:   nil,
			AndroidNotify: &AndroidNotify{
				Warn:    "1",            // 提醒类型： 1提示音；2震动；3指示灯，如果多个组合则对应编号组合如：12 标识提示音+震动
				Style:   0,              // 显示样式标识 0,"普通通知"，1,"BigTextStyle通知，点击后显示大段文字内容"，2,"BigPictureStyle，大图模式"，3,"横幅通知"，默认：0
				Content: []string{},     // 推送内容
				Sound:   p.AndroidSound, // 自定义声音
			}, // android通知消息, type=1, android
			IosNotify: &IosNotify{
				SubTitle:         "", // 副标题
				Sound:            p.ApnsSound,
				Badge:            "",             // 角标
				BadgeType:        0,              // badge类型, 1:绝对值 不能为负数，2增减(正数增加，负数减少)，减到0以下会设置为0
				CateGory:         p.ApnsCateGory, // apns的category字段，只有IOS8及以上系统才支持此参数推送
				SLIENT:           0,
				SlientPush:       0,  // 如果只携带content-available: 1,不携带任何badge，sound 和消息内容等参数， 则可以不打扰用户的情况下进行内容更新等操作即为“Silent Remote Notifications”
				ContentAvailable: 0,  // 将该键设为 1 则表示有新的可用内容。带上这个键值，意味着你的 App 在后台启动了或恢复运行了，application:didReceiveRemoteNotification:fetchCompletionHandler:被调用了
				MutableContent:   0,  // 需要在附加字段中配置相应参数
				AttachmentType:   0,  // ios富文本 0 无 ； 1 图片 ；2 视频 ；3 音频
				Attachment:       "", // ios富文本内容
			}, // ios通知消息, type=1, ios
			ExtrasMapList: pushMaps, // JSON格式 保留字段可参考下面附加参数示例
			Speed:         0,        // 每秒推送数量 只是趋势 默认0:不开启
		}
	}
	// PushTarget
	{
		model.PushTarget = PushTarget{
			Target:   2,   // 目标类型：1广播；2别名；3标签；4regid；5地理位置；6用户分群；9复杂地理位置推送
			Tags:     nil, // 客户端需要把 uid 设置为别名
			TagsType: "",  // target:3 => 标签组合方式：1并集；2交集；3补集(3暂不考虑)
			//Alias:     userIDList, // target:2 => 设置推送别名集合["alias1","alias2"]
			Rids:      nil, // target:4 => 设置推送Registration Id集合["id1","id2"]
			Block:     "",  // target:6 => 用户分群ID
			City:      "",  // target:5 => 推送地理位置 城市，地理位置推送时，city, province, country 必须有一个不为空
			Country:   "",  // target:5 => 推送地理位置 国家，地理位置推送时，city, province, country 必须有一个不为空
			Province:  "",  // target:5 => 推送地理位置 省份，地理位置推送时，city, province, country 必须有一个不为空
			PushAreas: nil, // target:9 时必传，复杂地理位置
		}
	}
	type response struct {
		Status int     `json:"status"`
		Error  *string `json:"error"`
	}
	// uid切片 最多1000个
	var uidList = make([][]string, (len(userIDList)+1000)/1000)
	for i := 0; i < len(userIDList); i += 1000 {
		end := i + 1000
		if end > len(userIDList) {
			end = len(userIDList)
		}
		uidList[i/1000] = userIDList[i:end]
	}
	for _, userIDList := range uidList {
		model.PushTarget.Alias = userIDList
		respBuf, err := p.client.Push(*model)
		if err != nil {
			logx.WithContext(ctx).Errorf("push error: %v", err)
			continue
		} else {
			var res response
			err = json.Unmarshal(respBuf, &res)
			if err == nil {
				if res.Status != 200 {
					logx.WithContext(ctx).Errorf("push error: %v", res.Error)
				}
			} else {
				logx.WithContext(ctx).Errorf("push error: %v", err)
			}
		}
		logx.WithContext(ctx).Infof("push success: %s", string(respBuf))
	}
	return "", nil
}
