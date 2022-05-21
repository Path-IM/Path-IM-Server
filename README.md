# Zero-IM-Server
基于 [Open-IM-Server](https://github.com/OpenIMSDK/Open-IM-Server) 实现的 IM 服务 

## 修改部分
### 服务注册发现
> 使用go-zero微服务框架 开发更方便 自带`链路追踪` `服务发现` `服务负载`

> 不依赖`mysql`所有业务逻辑均请求业务rpc服务 

### 新增超级大群功能
> 类似`QQ`群聊的`读扩散`模式  妈妈再也不用担心mongodb写入性能问题了

## 系统架构图
![system.svg](https://raw.githubusercontent.com/showurl/Zero-IM-Docs/main/images/20220517/Zero-IM-Server-System.svg)

## 业务架构图
![image1.svg](https://raw.githubusercontent.com/showurl/Zero-IM-Docs/main/images/20220517/Zero-IM-Server-Service.svg)

## 业务流程图
![flow.svg](https://raw.githubusercontent.com/showurl/Zero-IM-Docs/main/images/20220517/Zero-IM-Server-Flow.svg)

# Zero-IM-Server-Demo
> 使用 `Zero-IM-Server` 开发一个 `IM` 应用 
## 开发计划
- [x] 完成 Zero-IM-Server 的 TODO [第一天](https://github.com/showurl/Zero-IM-Server-Demo/tree/main/docs/day01)
- [x] 完成 用户模块 rpc 接口 编写 [第二天](https://github.com/showurl/Zero-IM-Server-Demo/tree/main/docs/day02)
- [x] 完成 用户关系模块 rpc 接口 编写 [第三天](https://github.com/showurl/Zero-IM-Server-Demo/tree/main/docs/day03/relation.md)
- [x] 完成 群聊模块 rpc 接口 编写 [第三天](https://github.com/showurl/Zero-IM-Server-Demo/tree/main/docs/day03/group.md)
- [x] 完成 用户模块 api 接口 编写 [第四天](https://github.com/showurl/Zero-IM-Server-Demo/tree/main/docs/day04)
- [x] 完成 用户关系模块 api 接口 编写 [第四天](https://github.com/showurl/Zero-IM-Server-Demo/tree/main/docs/day04)
- [x] 完成 群聊模块 api 接口 编写 [第四天](https://github.com/showurl/Zero-IM-Server-Demo/tree/main/docs/day04)

# Zero-IM-Client-Go
[Zero-IM-Client-Go](https://github.com/showurl/Zero-IM-Client-Go.git)
> 我们计划编写 `dart` sdk；由于时间问题，暂时放出 `golang` 客户端 测试代码；以供参考！

# 其他
## jaeger
![jaeger.png](https://raw.githubusercontent.com/showurl/Zero-IM-Docs/main/images/20220517/jaeger.png)
