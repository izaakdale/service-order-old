run_local: 
	PORT=8080 \
	HOST=localhost \
	AWS_REGION=eu-west-2 \
	TABLE_NAME=orders \
	TOPIC_ARN=arn:aws:sns:eu-west-2:000000000000:order-events \
	AWS_ENDPOINT=http://localhost:4566 \
	QUEUE_URL="http://localhost:4566/000000000000/order-queue" \
	go run .
run: 
	PORT=8080 \
	HOST=localhost \
	AWS_REGION=eu-west-2 \
	TABLE_NAME=orders \
	TOPIC_ARN=arn:aws:sns:eu-west-2:000000000000:order-notify \
	go run .


PROTO_DIR=schema/order

gproto:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative ${PROTO_DIR}/*.proto \
	 --go-grpc_out=. --go-grpc_opt=paths=source_relative ${PROTO_DIR}/*.proto