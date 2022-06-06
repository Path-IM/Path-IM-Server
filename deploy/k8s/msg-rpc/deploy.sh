#!/usr/bin/env bash
tag=`date +%Y%m%d%H%M%S`
echo "tag: $tag"
rm -rf ./bin
GOOS=linux GOARCH=amd64 go build -o bin ../../../app/msg/cmd/rpc/chat.go || exit 1
docker build --platform linux/amd64 . -t ccr.ccs.tencentyun.com/zeroim/msg-rpc:${tag} || exit 1
docker push ccr.ccs.tencentyun.com/zeroim/msg-rpc:${tag} || exit 1
rm -rf ./deployment.yaml
cp deployment.tmp.yaml deployment.yaml
sed -i "" "s#TMP_IMAGE#ccr.ccs.tencentyun.com/zeroim/msg-rpc:${tag}#g" deployment.yaml
