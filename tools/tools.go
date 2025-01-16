// To install the following tools at the version used by this repo run:
// $ make tools
// or
// $ go generate -tags tools tools/tools.go

package tools

//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate go install github.com/bufbuild/buf/cmd/buf
//go:generate go install github.com/kyleconroy/sqlc/cmd/sqlc

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/kyleconroy/sqlc/cmd/sqlc"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
