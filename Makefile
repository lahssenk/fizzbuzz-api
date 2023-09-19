gobench:
	go test -bench=. --benchmem ./...

PROTOS = $(shell find . -iname '*.proto')


.PHONY: protoc 
protoc:
	protoc -I=./proto --go_out=./pkg/protogen --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/protogen --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=./pkg/protogen --grpc-gateway_opt paths=source_relative \
	--openapi_out=./openapi_v3 --openapi_opt=enum_type=string,default_response=true \
	./proto/fizzbuzz/v1/*.proto

govet:
	go vet ./...

gotest: govet
	go test ./... 
	
gorun: gotest
	API_KEY=testkey SERVER_HOST=localhost SERVER_PORT=8080 ADMIN_PORT=8081 go run ./cmd/api/...

build:
	docker build -t fizzbuzz .

ci: build
	docker run --rm -it -p 8080:8080 -p 8081:8081 fizzbuzz
