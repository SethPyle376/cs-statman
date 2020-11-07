build: copy-protos statman-processor statman statman-api statman-ui

copy-protos:
	cp pkg/protos/*.proto api/protos/

gen-protos:
	if [ ! -d "pkg/csproto" ]; then mkdir pkg/csproto; fi && protoc -I ./pkg/protos ./pkg/protos/*.proto --go-grpc_out=./pkg/csproto --go_out=./pkg/csproto

statman-processor:
	docker build -t statman-processor -f ./processor/Dockerfile .

statman:
	docker build -t statman -f ./statman/Dockerfile .

statman-api:
	docker build -t statman-api -f ./api/Dockerfile ./api

statman-ui:
	docker build -t statman-ui -f ./statman-ui/Dockerfile ./statman-ui