#!/usr/bin/env bash
#export PATH=$PATH:$HOME/sdk/flutter/2.5.3/.pub-cache/bin
#flutter pub global activate protoc_plugin
#protoc param.proto --dart_out=./dart
rm -rf ./*.dart
protoc -I . *.proto --dart_out=.
tar -zcvf ~/Desktop/dart.pb.tar.gz ./*