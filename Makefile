gen:
	protoc -I ./internal/proto/src \
    --go_out ./internal/proto/gen --go_opt paths=source_relative  \
    --go-grpc_out ./internal/proto/gen --go-grpc_opt paths=source_relative \
    auth.proto

install_protoc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

srv:
	go build ./cmd/server/main.go

cert:
	openssl req -x509 -newkey rsa:4096 -keyout ./data/cert/priv.pem -out ./data/cert/cert.pem -passout pass:6879hjkw%4 -sha256 -days 365
	openssl rsa -in ./data/cert/priv.pem -out ./data/cert/priv.pem -passin pass:6879hjkw%4