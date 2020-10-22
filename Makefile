build: gen-protos

gen-protos:
	if [ ! -d "csproto" ]; then mkdir csproto; fi && protoc -I ./protos ./protos/*.proto --go-grpc_out=./csproto