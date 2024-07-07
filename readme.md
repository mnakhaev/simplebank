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

How to generate swagger documentation:
- Add `--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true` to `protoc` command. 
`allow_merge=true` is needed to have single JSON file instead of multiple.
- Add next section to automate documentation info:
```swagger codegen
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  info: {
    title: "Simple Bank API";
    version: "1.0";
    contact: {
      name: "gRPC-Gateway project";
      url: "https://github.com/mnakhaev";
      email: "maga2894@gmail.com";
    };
  }
```
- Fix the imports (check data below)
- Download swagger UI repo and copy some files: `cp -r swagger-ui/dist/* doc/swagger`
- Edit copied files (point on correct swagger file)
- Add new handler for Swagger

If import in proto file is not successful, then:
1) Copy repository which is imported. You will need it to copy some missing file(s).
2) Create same structure for needed file in your local repository - `mkdir -p proto/protoc-gen-openapiv2/options`
3) Copy all proto files from downloaded repo to local folder - `cp grpc-gateway/protoc-gen-openapiv2/options/*.proto proto/protoc-gen-openapiv2/options`

[statik](https://github.com/rakyll/statik) is used to embed static files into Golang binary.
Usage in code:
```go
statikFS, err := fs.New()
if err != nil {
    log.Fatal("cannot create statik fs:", err)
}

swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
mux.Handle("/swagger/", swaggerHandler)
```
Note: now it's outdated, use https://go.dev/doc/go1.16#library-embed instead.


How to add new gRPC API:
1) Create new proto file in /proto directory
2) Add import and RPC description in `service_simple_bank.proto`
3) Run `make proto`