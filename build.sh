SRC_DIR=./finite
DST_DIR=./finite

protoc -I=$SRC_DIR \
  --go_out=plugins=grpc:$DST_DIR \
  $SRC_DIR/event.proto

