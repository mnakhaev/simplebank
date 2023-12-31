Create migrations:
- `migrate create -ext sql -dir db/migration -seq init_schema`

Create Postgres DB:
- `createdb --username=root --owner=root simple_bank`

Init SQLC:
- `sqlc init`

Generate code using SQLC:
- `sqlc generate`