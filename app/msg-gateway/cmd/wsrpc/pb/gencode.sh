#!/bin/bash
goctl rpc protoc -I=./ -I=../../../../msg/cmd/rpc/pb msg-gateway.proto -v --go_out=.. --go-grpc_out=..  --zrpc_out=.. --style=goZero --home ../../../../../goctl/home
#protoc -I=../../../../msg/cmd/rpc/pb msg-gateway.proto --proto_path ./ --go_out .. --go-grpc_out ..
