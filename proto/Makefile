PROTO_ROOT=.

all: install generate_protos

generate_protos:
	protoc -I${PROTO_ROOT}/product-catalog-service -I${PROTO_ROOT}/third_party -I${PROTO_ROOT}/third_party/googleapis --openapiv2_out=logtostderr=true:./product-catalog-service --openapiv2_opt=generate_unbound_methods=true ${PROTO_ROOT}/product-catalog-service/product.proto
	protoc -I${PROTO_ROOT}/product-categories-service -I${PROTO_ROOT}/third_party -I${PROTO_ROOT}/third_party/googleapis --openapiv2_out=logtostderr=true:./product-categories-service --openapiv2_opt=generate_unbound_methods=true ${PROTO_ROOT}/product-categories-service/category.proto

install:
	go mod download
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc \
        google.golang.org/protobuf/cmd/protoc-gen-go