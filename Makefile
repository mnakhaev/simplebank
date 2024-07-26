DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrate-up:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrate-up-aws:
	migrate -path db/migration -database "postgresql://root:***@simple-bank.cr2iqsqagfht.eu-north-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migrate-up-last:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migrate-down:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrate-down-last:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new-migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

rundb: postgres createdb migrate-up

sqlc:
	sqlc generate

server:
	go run cmd/main.go

test:
	go test -v -cover -short ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/mnakhaev/simplebank/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb \
 	--go-grpc_opt=paths=source_relative --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
 	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
 	proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl
	# package pb -> service SimpleBank -> call CreateUser

redis:
	docker run --name redis --network bank-network -p 6379:6379 -d redis:7.4-rc2-alpine

.PHONY: postgres createdb dropdb migrate-up migrate-down sqlc mock migrate-up-last migrate-down-last new-migration rundb db_docs db_schema proto evans redis