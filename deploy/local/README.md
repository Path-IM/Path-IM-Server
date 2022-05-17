# 开发环境搭建

## goctl

```shell
cd goctl/goctl
go install .
```

## 开源组件依赖

### docker

#### etcd

##### arm64

```shell
docker run --name etcd -p 2379:2379 -p 2380:2380 -d showurl/etcd-arm /usr/local/bin/etcd -advertise-client-urls http://0.0.0.0:2379 -listen-client-urls http://0.0.0.0:2379
```

#### etcd-keeper

```shell
docker run --name etcd-keeper -d --link etcd:etcd -p 8080:8080 showurl/etcdkeeper:v0.7.6 sh -c "/etcdkeeper/etcdkeeper || /opt/etcdkeeper/etcdkeeper.bin -h \$HOST -p \$PORT"
open http://localhost:8080
# 输入 etcd:2379
```

#### es

```shell
docker run --name elasticsearch -d \
-p 9200:9200 \
-p 9300:9300 \
-e "discovery.type=single-node" \
-e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
-e "TZ=Asia/Shanghai" \
--user root \
--privileged \
docker.elastic.co/elasticsearch/elasticsearch:7.13.4
```

#### kibana

```shell
docker run --name kibana -d \
-p 5601:5601 \
--link elasticsearch:elasticsearch \
-e "elasticsearch.hosts=http://elasticsearch:9200" \
-e "TZ=Asia/Shanghai" \
-e "I18N_LOCALE=zh-CN" \
--privileged \
docker.elastic.co/kibana/kibana:7.13.4
```

#### jaeger

```shell
docker run --name jaeger -d \
-p 5775:5775/udp \
-p 6831:6831/udp \
-p 6832:6832/udp \
-p 5778:5778 \
-p 16686:16686 \
-p 14268:14268 \
-p 9411:9411 \
-e SPAN_STORAGE_TYPE=elasticsearch \
-e ES_SERVER_URLS=http://elasticsearch:9200 \
-e LOG_LEVEL=debug \
--link elasticsearch:elasticsearch \
jaegertracing/all-in-one:latest 
open http://localhost:16686/
```

#### prometheus
```shell
docker run --name prometheus -d \
-e "TZ=Asia/Shanghai" \
-p 9090:9090 \
-v $(pwd)/deploy/local/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml \
--user root \
--privileged \
prom/prometheus:v2.28.1 \
--config.file=/etc/prometheus/prometheus.yml \
--storage.tsdb.path=/prometheus
```
#### grafana
```shell
docker run --name grafana -d \
-e "TZ=Asia/Shanghai" \
-p 3000:3000 \
--user root \
--privileged \
grafana/grafana:8.0.6
```

#### mongodb
```shell
docker run --name mongodb -d \
-p 27017:27017 \
--user root \
--privileged \
mongo:latest
```
#### redis 
```shell
docker run --name redis -d \
-p 6379:6379 \
--user root \
--privileged \
crysislinux/rebloom:2.2.9 \
redis-server --requirepass 123456 --maxmemory 1gb --maxmemory-policy allkeys-lru --loadmodule /usr/lib/redis/modules/redisbloom.so
```
### docker-compose
#### zookeeper kafka kafka-ui 
```shell
cd deploy/local/kafka
docker-compose up -d
```

# 运行服务
## 环境变量
```shell
GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn
```
## 以下服务依次运行
### im-user rpc
### msg-callback rpc
### conversation rpc
### auth rpc
### msg-push rpc
### msg-transfer history
### msg rpc