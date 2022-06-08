- docker-compose.yaml x86架构
> wget https://github.91chi.fun//https://github.com//docker/compose/releases/download/v2.5.1/docker-compose-linux-x86_64
> chmod +x docker-compose-linux-x86_64 && mv docker-compose-linux-x86_64 /usr/bin/docker-compose
# 依赖
>dependencies/docker-compose.yaml
>替换ip地址[10.1.3.12]为内网/公网ip

```shell
cd dependencies
docker-compose up -d
```
# kafka topic 提前创建：
- im_msg
- im_msg_push_single
- im_msg_push_group

# 服务
> *.yaml
> 替换ip地址[10.1.3.12]为内网/公网ip
```shell
cd ..
docker-compose up -d
```