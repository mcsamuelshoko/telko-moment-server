start:
	go run ./cmd/server/main.go

tidy:
	go mod tidy

codegen:
	oapi-codegen \
	-generate "fiber,types,strict-server,spec" \
	-package=api -o api/api.gen.go oapi_codegen.yml