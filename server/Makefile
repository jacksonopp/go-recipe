build-linux:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/server-linux-amd64 cmd/server/main.go

build-mac:
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/server-mac-amd64 cmd/server/main.go

run-mac:
	@./bin/server-mac-amd64