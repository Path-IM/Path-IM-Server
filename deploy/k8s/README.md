# 创建命名空间
```shell
kubectl create namespace pathim
```
# 首先要创建role
```shell
kubectl apply -f roles.yaml
```
# 镜像认证 registry

# pathim-configmap
```shell
kubectl apply -f pathim-configmap.yaml
```
# 依赖
## cassandra
```shell
cd cassandra
kubectl apply -f . 
```
## redis
```shell
cd redis
kubectl apply -f .
```
# 各服务
## imuser-rpc
```shell
cd imuser-rpc
bash deploy.sh
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
```
## msgcallback-rpc
```shell
cd msgcallback-rpc
bash deploy.sh
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
```
## msg-rpc
```shell
cd msg-rpc
bash deploy.sh
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
```
## msggateway-wsrpc
```shell
cd msggateway-wsrpc
bash deploy.sh
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
```
## msgpush-rpc
```shell
cd msgpush-rpc
bash deploy.sh
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
```
## transfer-history-cassandra
```shell
cd transfer-history-cassandra
bash deploy.sh
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
```