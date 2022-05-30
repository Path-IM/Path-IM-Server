#!/bin/bash
goctl rpc protoc -I=./ chat.proto -v --go_out=.. --go-grpc_out=..  --zrpc_out=.. --style=goZero --home ../../../../../goctl/home
goctl rpc protoc -I=./ ws.proto -v --go_out=.. --go-grpc_out=..  --zrpc_out=.. --style=goZero --home ../../../../../goctl/home
sed -i "" 's#pb "github.com/Path-IM/Path-IM-Server/app/msg/cmd/rpc/pb"##g' chat.pb.go
sed -i "" 's#pb.##g' chat.pb.go
mv ../github.com/Path-IM/Path-IM-Server-Demo/app/msg/cmd/rpc/pb/ws.pb.go ./
rm -rf ../github.com
rm -rf ../example
rm -rf ../etc/ws.yaml
rm -rf ../internal/server/exampleServer.go
rm -rf ../ws.go
