# protoc-gen-go-grpc-client

The protoc-gen-go-grpc-client plugin generates a
Go client for a gRPC API interface while allowing the seamless transfer of supported metadata values from the context.

## Why this plugin?
In gRPC, the context carries important information that can be used for various purposes,
such as authentication tokens, request tracing, and user-specific data. However,
by default, the context values are not directly accessible within the gRPC client. 
The protoc-gen-go-grpc-client plugin addresses this limitation by enabling the transfer of these
context values to the gRPC client via gRPC metadata API.


## Installation

```bash
go install github.com/rmnkmr/lsp/plugins/protoc-gen-go-grpc-client
```

## Usage

```bash
protoc --go-grpc-client_out=paths=source_relative:. path/to/file.proto
```