#!/bin/sh -e

INC=$(go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway)
INC_OPENAPI=$(go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway/v2)
ARGS="-I${INC}"
ARGS_OPENAPI="-I${INC_OPENAPI}"

protoc $ARGS/third_party/googleapis $ARGS_OPENAPI -Iuser --grpc-gateway_out ./user --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true  --swagger_out=allow_merge=true,merge_file_name=api:. \
 --go_out ./user --go_opt paths=source_relative --go-grpc_out ./user --go-grpc_opt paths=source_relative ./user/*.proto

protoc $ARGS/third_party/googleapis $ARGS_OPENAPI -Icompany --grpc-gateway_out ./company --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true  --swagger_out=allow_merge=true,merge_file_name=api:. \
  --go_out ./company --go_opt paths=source_relative --go-grpc_out ./company --go-grpc_opt paths=source_relative ./company/*.proto