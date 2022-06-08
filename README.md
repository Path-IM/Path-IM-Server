# Path-IM-Server
使用go-zero框架开发的IM服务器。 有高度定制IM需求的开发者，可以使用这个项目。
> 普通开发者可以在[演示项目](https://github.com/Path-IM/Path-IM-Server-Demo) 基础上进行开发。

> [文档入口](https://doc.pathim.cn)
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
![image1.svg](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220608/Path-IM-Server%E4%B8%9A%E5%8A%A1%E6%9E%B6%E6%9E%84%E5%9B%BE.svg)

## 业务流程图
![flow.svg](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220608/Path-IM-Server%E4%B8%9A%E5%8A%A1%E6%B5%81%E7%A8%8B%E5%9B%BE.svg)

# 部署运行
## docker-compose
> 目录：deploy/local/pathim-docker

- ！！！一定要先替换内网地址 目录下文件全部替换 替换ip地址`10.1.3.12`为`内网/公网ip`
- ！！！一定要先替换内网地址 目录下文件全部替换 替换ip地址`10.1.3.12`为`内网/公网ip`
- ！！！一定要先替换内网地址 目录下文件全部替换 替换ip地址`10.1.3.12`为`内网/公网ip`
### 如何按照docker-compose
#### linux
```shell
wget https://github.91chi.fun//https://github.com//docker/compose/releases/download/v2.5.1/docker-compose-linux-x86_64
chmod +x docker-compose-linux-x86_64 && mv docker-compose-linux-x86_64 /usr/bin/docker-compose
```
### 依赖
```shell
cd deploy/local/pathim-docker/dependencies
docker-compose up -d
```
> 打开`内网/公网ip`:8081 进入kafka-ui 主动创建以下topic

- im_msg
- im_msg_push_single
- im_msg_push_group
[img.png](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220608/kafkaui.png)
### Path-IM-Server各服务
```shell
cd deploy/local/pathim-docker
docker-compose up -d
```
### 服务运行情况
[img.png](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220608/docker-compose.png)

### 
## 源码部署
### 编译命令
```shell
go build -o bin .
```
### Dockerfile
```dockerfile
FROM showurl/zerobase
WORKDIR /app
COPY ./bin /app/zeroservice
RUN chmod +x /app/zeroservice && mkdir /app/etc
CMD ["/app/zeroservice"]
```
### docker容器中运行
```shell
docker run -v ./xxx.yaml:/app/etc/xxx.yaml your-image:tag
```

# 其他
## jaeger
![jaeger.png](https://raw.githubusercontent.com/Path-IM/Path-IM-Docs/main/images/20220517/jaeger.png)
## 如果选择使用mongo或cassandra
- 1、msg-transfer中只部署`history-cassandra`服务，或`history-mongo`服务
- 2、msg-rpc配置文件`HistoryDBType`设置为`cassandra`/`mongo`

> 本项目设计思路源自[Open-IM-Server](https://github.com/OpenIMSDK/Open-IM-Server)