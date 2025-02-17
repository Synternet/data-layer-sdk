package main

//go:generate protoc -I../../proto -Iproto --go_out=types --go_opt=paths=source_relative,Mgoogle/protobuf/empty.proto=google.golang.org/protobuf/types/known/emptypb --go-grpc_out=types --go-grpc_opt=paths=source_relative,Mgoogle/protobuf/empty.proto=google.golang.org/protobuf/types/known/emptypb proto/example/v1/example.proto
