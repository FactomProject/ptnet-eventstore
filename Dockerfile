FROM golang:1.12

WORKDIR /app

# install protobuf/grpc tools
RUN \
      go get -u github.com/golang/protobuf/protoc-gen-go && \
      go get -u google.golang.org/grpc

ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV PFLOWPATH=/app/examples
EXPOSE 50051
ENTRYPOINT ./entry.sh
