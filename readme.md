Create migrations:
- `migrate create -ext sql -dir db/migration -seq init_schema`

Create Postgres DB:
- `createdb --username=root --owner=root simple_bank`

Init SQLC:
- `sqlc init`

Generate code using SQLC:
- `sqlc generate`

Run server in production mode (override default DB source):
- `docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres12:5432/simple_bank?sslmode=disable" simplebank:latest`