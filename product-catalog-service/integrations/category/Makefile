PROTO_ROOT=../../../proto

all: clean
	protoc -I${PROTO_ROOT}/product-categories-service \
 		   -I${PROTO_ROOT}/third_party \
 		   -I${PROTO_ROOT}/third_party/googleapis \
 		   --go_out . \
 		   --go-grpc_out . \
           --go-grpc_opt=require_unimplemented_servers=false \
 		   ${PROTO_ROOT}/product-categories-service/category.proto

clean:
	rm -rf *.pb.go