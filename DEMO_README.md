# DEMO
## fork/clone 仓库
```shell
git clone https://github.com/showurl/Zero-IM-Server.git -b main --depth=1
```
## 搜索 todo 的代码 实现自己的逻辑
![todo.png](https://raw.githubusercontent.com/showurl/Zero-IM-Docs/main/images/20220517/todo.svg)

## 服务启动运行
> 注意⚠️ 要有顺序的运行

### imuser-rpc
> 你自己业务的逻辑 包括 `验证token`， `获取群成员`， `是否被拉黑`， `是否是好友`， `获取会话推送消息设置`。

### msg-callback-rpc
> 你自己业务的逻辑 消息回调

> 可通过 msg-rpc 服务进行配置 `disable`

### msg-rpc
> 用户发送消息的处理

#### 配置
```yaml
MessageVerify:
  FriendVerify: true # 只有好友可以发送消息

Callback:
  CallbackWordFilter:
    Enable: true # 开启关键词过滤
    ContinueOnError: true # 开启关键词过滤时，如果出错，是否继续发送
  CallbackAtAllInSuperGroup:
    Enable: true # 超级大群中，是否允许@所有人
    ContinueOnError: true # 超级大群中，@所有人时，如果出错，是否继续发送通知
  CallbackBeforeSendGroupMsg:
    Enable: true # 开启群消息发送前回调
    ContinueOnError: true # 开启群消息发送前回调时，如果出错，是否继续发送
  CallbackAfterSendGroupMsg:
    Enable: true # 开启群消息发送后回调
    ContinueOnError: true # 无意义
  CallbackBeforeSendSuperGroupMsg:
    Enable: true # 开启超级大群消息发送前回调
    ContinueOnError: true # 开启超级大群消息发送前回调时，如果出错，是否继续发送
  CallbackAfterSendSuperGroupMsg:
    Enable: true # 开启超级大群消息发送后回调
    ContinueOnError: true # 无意义
  CallbackBeforeSendSingleMsg:
    Enable: true # 开启私聊消息发送前回调
    ContinueOnError: true # 开启私聊消息发送前回调时，如果出错，是否继续发送
  CallbackAfterSendSingleMsg:
    Enable: true # 开启私聊消息发送后回调
    ContinueOnError: true # 无意义
```

### msg-push-rpc
> 消息在线和离线推送

### msg-transfer-history
> 离线消息数据传输 写入 `mongodb`、 写 `seq`

### msg-gateway
> websocket server
#### 配置
```yaml
Websocket:
  MaxConnNum: 10000 # 单pod最大连接用户数量
  TimeOut: 10 # HandshakeTimeout
  MaxMsgLen: 4096 # 指定I/O缓冲区大小（以字节为单位）如果缓冲区大小为零，则使用HTTP服务器分配的缓冲区。I/O缓冲区的大小不限制可以发送或接收的消息的大小。
```
> rpc: 获取在线用户