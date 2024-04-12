PROTO_ROOT=../../../proto

all: clean install
	protoc -I${PROTO_ROOT}/product-categories-service \
 		   -I${PROTO_ROOT}/third_party \
 		   -I${PROTO_ROOT}/third_party/googleapis \
			--go_out . \
			--go-grpc_out . \
			--go-grpc_opt=require_unimplemented_servers=false \
			--grpc-gateway_out=logtostderr=true:. \
			--govalidators_out=. \
 		   ${PROTO_ROOT}/product-categories-service/category.proto \

clean:
	rm -f category.pb.go category.pb.gw.go category.validator.pb.go category_grpc.pb.go

install:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
