# Path-IM-Server
使用go-zero框架开发的IM服务器。 有高度定制IM需求的开发者，可以使用这个项目。
> 普通开发者可以在[演示项目](https://github.com/Path-IM/Path-IM-Server-Demo) 基础上进行开发。

> [文档入口](https://pathim.msypy.xyz)
## 优势
- 使用go-zero微服务框架 开发更方便 自带`链路追踪`,`p2c服务负载均衡`,`熔断限流`,`自适应降载`等功能
- 不依赖`mysql`所有业务逻辑均请求你自己的业务rpc接口 你只需实现rpc接口即可 
- 可以使用 `cassandra` 来替代 `mongodb`
- 类似`QQ`群聊的`读扩散`模式  妈妈再也不用担心`mongodb`/`cassandra`写入性能问题了
- 使用dart开发[sdk](https://github.com/Path-IM/Path-IM-Core-Flutter), 使用flutter做客户端, 直接生成5端代码

## 开源组件依赖
- mongodb or cassandra (离线消息存储 个人推荐cassandra)
- kafka (消息队列)
- redis (存储seq)
- ~~etcd~~ (不依赖etcd)
- ~~mysql~~ (不依赖mysql)

# rancher 查看服务运行情况
- 地址：[rancher](https://42.194.149.177:1443)
- 用户名：guest
- 密码：guest

## 业务架构图
![image1.svg](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220604/Path-IM-Server%20%E4%B8%9A%E5%8A%A1%E6%9E%B6%E6%9E%84%E5%9B%BE.svg)

## 业务流程图
![flow.svg](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220604/Path-IM-Server%E6%B6%88%E6%81%AF%E6%B5%81%E8%BD%AC%E5%9B%BE.svg)

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

# 其他
## jaeger
![jaeger.png](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220517/jaeger.png)
## 如果选择使用mongo或cassandra
- 1、msg-transfer中只部署`history-cassandra`服务，或`history-mongo`服务
- 2、msg-rpc配置文件`HistoryDBType`设置为`cassandra`/`mongo`

> 本项目设计思路源自[Open-IM-Server](https://github.com/OpenIMSDK/Open-IM-Server)