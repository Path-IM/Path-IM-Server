# Zero-IM-Server
基于 [Open-IM-Server](https://github.com/OpenIMSDK/Open-IM-Server) 实现的 IM 服务 

## 修改部分
### 服务注册发现
> 使用go-zero微服务框架 开发更方便 自带 `链路追踪` `服务发现` `服务负载`

![jaeger.png](https://public.msypy.xyz/images/Zero-IM-Server/477E4F088A59947E34D949A584E39A62.jpg)

> 不依赖 `mysql` 所有业务逻辑均请求业务rpc服务 
