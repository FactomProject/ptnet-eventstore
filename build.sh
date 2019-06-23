SRC_DIR=./finite
DST_DIR=./finite

# generate protobuf code
protoc -I=$SRC_DIR \
  --go_out=plugins=grpc:$DST_DIR \
  $SRC_DIR/event.proto

# build the binary
go build -o ptneteventstore ./main.go 
