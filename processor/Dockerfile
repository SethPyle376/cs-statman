FROM golang:1.14-alpine as Builder

ADD . /go/src/github.com/sethpyle376/cs-statman

WORKDIR /go/src/github.com/sethpyle376/cs-statman/processor

RUN apk add --no-cache ca-certificates protobuf git

RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN cp /go/bin/protoc-gen-go /usr/bin
RUN cp /go/bin/protoc-gen-go-grpc /usr/bin

if [ ! -d "/go/src/github.com/sethpyle376/cs-statman/pkg/csproto" ]; then mkdir /go/src/github.com/sethpyle376/cs-statman/pkg/csproto; fi;
RUN protoc -I /go/src/github.com/sethpyle376/cs-statman/pkg/protos /go/src/github.com/sethpyle376/cs-statman/pkg/protos/*.proto --go-grpc_out=/go/src/github.com/sethpyle376/cs-statman/pkg/csproto --go_out=/go/src/github.com/sethpyle376/cs-statman/pkg/csproto

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o processor

FROM alpine:3.12

RUN apk add --no-cache ca-certificates

COPY --from=Builder /go/src/github.com/sethpyle376/cs-statman/processor/processor /processor
RUN mkdir /temp

CMD ["/processor"]