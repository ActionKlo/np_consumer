include .env

#migrations
gooseUp:
	goose -dir internal/db/migrations postgres postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB} up
gooseReset:
	goose -dir internal/db/migrations postgres postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB} reset
gooseReUp: gooseReset gooseUp

startConsumer:
	go run ./cmd/consumer/

#grpc
grpc:
	protoc ./proto/*.proto \
		--go_out=./internal/api \
		--go_opt=paths=source_relative \
		--go-grpc_out=./internal/api \
		--go-grpc_opt=paths=source_relative

new: gooseReUp startConsumer
