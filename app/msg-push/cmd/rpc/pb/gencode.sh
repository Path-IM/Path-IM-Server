#!/bin/bash
goctl rpc protoc -I=./ -I=../../../../msg/cmd/rpc/pb msgpush.proto -v --go_out=.. --go-grpc_out=..  --zrpc_out=.. --style=goZero --home ../../../../../goctl/home

