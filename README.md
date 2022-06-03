# Path-IM-Server
使用go-zero框架开发的IM服务器。 有高度定制IM需求的开发者，可以使用这个项目。
> 普通开发者可以在[演示项目](https://github.com/Path-IM/Path-IM-Server) 基础上进行开发。

## 优势
- 使用go-zero微服务框架 开发更方便 自带`链路追踪`,`p2c服务负载均衡`,`熔断限流`,`自适应降载`等功能
- 不依赖`mysql`所有业务逻辑均请求你自己的业务rpc接口 你只需实现rpc接口即可 
- 可以使用 `cassandra` 来替代 `mongodb`
- 类似`QQ`群聊的`读扩散`模式  妈妈再也不用担心`mongodb`/`cassandra`写入性能问题了
## 开源组件依赖
- mongodb or cassandra (离线消息存储 个人推荐cassandra)
- kafka (消息队列)
- redis (存储seq)
- ~~etcd~~ (不依赖etcd)
- ~~mysql~~ (不依赖mysql)

## 系统架构图
![system.svg](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220517/Path-IM-Server-System.svg)

## 业务架构图
![image1.svg](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220517/Path-IM-Server-Service.svg)

## 业务流程图
![flow.svg](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220517/Path-IM-Server-Flow.svg)

# 演示项目:Path-IM-Server (最好的文档就是demo)
> 使用 `Path-IM-Server` 开发一个 `IM` 应用 
## 开发计划
- [x] 完成 Path-IM-Server 的 TODO [第一天](https://github.com/Path-IM/Path-IM-Server/tree/main/docs/day01)
- [x] 完成 用户模块 rpc 接口 编写 [第二天](https://github.com/Path-IM/Path-IM-Server/tree/main/docs/day02)
- [x] 完成 用户关系模块 rpc 接口 编写 [第三天](https://github.com/Path-IM/Path-IM-Server/tree/main/docs/day03/relation.md)
- [x] 完成 群聊模块 rpc 接口 编写 [第三天](https://github.com/Path-IM/Path-IM-Server/tree/main/docs/day03/group.md)
- [x] 完成 用户模块 api 接口 编写 [第四天](https://github.com/Path-IM/Path-IM-Server/tree/main/docs/day04)
- [x] 完成 用户关系模块 api 接口 编写 [第四天](https://github.com/Path-IM/Path-IM-Server/tree/main/docs/day04)
- [x] 完成 群聊模块 api 接口 编写 [第四天](https://github.com/Path-IM/Path-IM-Server/tree/main/docs/day04)
- [x] 完成 k8s 部署 [第五天](https://github.com/Path-IM/Path-IM-Server/tree/main/deploy/k8s)
- [x] 完成 api 文档 编写 [第六天](https://github.com/Path-IM/Path-IM-Server/tree/main/docs/api.md)
- [x] 完成 消息持久化存储 文档 编写 [第十天](https://github.com/Path-IM/Path-IM-Server/tree/main/docs/persistent.md)
- [x] 支持 cassandra 离线消息存储 [第十二天](https://github.com/Path-IM/Path-IM-Server/tree/main/docs/cassandra.md)
- [x] go-zero periodlimit 用户发送消息限流 [第十二天](https://github.com/Path-IM/Path-IM-Server/tree/main/docs/periodlimit.md)

# Path-IM-Client-Go
[Path-IM-Client-Go](https://github.com/Path-IM/Path-IM-Client-Go.git)
> 我们计划编写 `dart` sdk；由于时间问题，暂时放出 `golang` 客户端 测试代码；以供参考！

# 其他
## jaeger
![jaeger.png](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220517/jaeger.png)
## 如果选择使用mongo或cassandra
- 1、msg-transfer中只部署`history-cassandra`服务，或`history-mongo`服务
- 2、msg-rpc配置文件`HistoryDBType`设置为`cassandra`/`mongo`

> 本项目设计思路源自[Open-IM-Server](https://github.com/OpenIMSDK/Open-IM-Server)