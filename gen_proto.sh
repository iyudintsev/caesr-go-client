protoc -I./proto \
       --go_opt=module=github.com/iyudintsev/caesr-go-client \
       --go_out=./proto \
       --go-grpc_opt=module=github.com/iyudintsev/caesr-go-client \
       --go-grpc_out=./proto proto/caesr.proto \
       --experimental_allow_proto3_optional
