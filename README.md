# Zero-IM-Server
基于 [Open-IM-Server](https://github.com/OpenIMSDK/Open-IM-Server) 实现的 IM 服务 

## 修改部分
### 服务注册发现
> 使用go-zero微服务框架 开发更方便 自带 `链路追踪` `服务发现` `服务负载`

![jaeger.png](https://public.msypy.xyz/images/Zero-IM-Server/477E4F088A59947E34D949A584E39A62.jpg)

> 不依赖 `mysql` 所有业务逻辑均请求业务rpc服务 

## 系统架构图
![system.svg](https://public.msypy.xyz/images/Zero-IM-Server/Zero-IM-Server%20%20%E7%B3%BB%E7%BB%9F%E6%9E%B6%E6%9E%84%E5%9B%BE.svg)

## 业务架构图
![image1.svg](https://public.msypy.xyz/images/Zero-IM-Server/Zero-IM-Server%20%E4%B8%9A%E5%8A%A1%E6%9E%B6%E6%9E%84%E5%9B%BE.svg)

## 业务流程图
![flow.svg](https://public.msypy.xyz/images/Zero-IM-Server/Zero-IM-Server-Flow.svg)

# Zero-IM-Server-Demo
[Zero-IM-Server-Demo](DEMO_README.md)

# Zero-IM-Client-Go
[Zero-IM-Client-Go](https://github.com/showurl/Zero-IM-Client-Go.git)
> 我们计划编写 `dart` sdk；由于时间问题，暂时放出 `golang` 客户端 测试代码；以供参考！
