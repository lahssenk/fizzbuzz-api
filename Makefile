gobench:
	go test -bench=. --benchmem ./...

govet:
	go vet ./...

gotest: govet
	go test ./... 
	
gorun: gotest
	API_KEY=testkey SERVER_HOST=localhost SERVER_PORT=8080 ADMIN_PORT=8081 go run ./cmd/main.go

build:
	docker build -t fizzbuzz .

ci: build
	docker run --rm -it -p 8080:8080 -p 8081:8081 fizzbuzz
