FROM golang:1.14-alpine as Builder

ADD . /go/src/github.com/sethpyle376/cs-statman

WORKDIR /go/src/github.com/sethpyle376/cs-statman

RUN apk add --no-cache ca-certificates protobuf

RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN cp /go/bin/protoc-gen-go /usr/bin
RUN cp /go/bin/protoc-gen-go-grpc /usr/bin

RUN if [ ! -d "/go/src/github.com/sethpyle376/cs-statman/pkg/csproto" ]; then mkdir /go/src/github.com/sethpyle376/cs-statman/pkg/csproto; fi
RUN protoc -I /go/src/github.com/sethpyle376/cs-statman/pkg/protos /go/src/github.com/sethpyle376/cs-statman/pkg/protos/*.proto --go_out=./pkg/csproto --go_opt=paths=source_relative \
    --go-grpc_out=./pkg/csproto --go-grpc_opt=paths=source_relative

WORKDIR /go/src/github.com/sethpyle376/cs-statman/statman

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o statman

FROM alpine:3.12

RUN apk add --no-cache ca-certificates

COPY --from=Builder /go/src/github.com/sethpyle376/cs-statman/statman/statman /statman

CMD ["/statman"]