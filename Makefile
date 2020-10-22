build: gen-protos

gen-protos:
	if [ ! -d "pkg/csproto" ]; then mkdir pkg/csproto; fi && protoc -I ./pkg/protos ./pkg/protos/*.proto --go-grpc_out=./pkg/csproto --go_out=./pkg/csproto