all: build

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

build: go.sum
	@swag init
	@go mod tidy
	@go fmt ./...
	@go build  -o build/uploader .

.PHONY: all build build-linux build-window


.PHONY: release
release: CGO_ENABLED=0
release:
	GOOS=windows GOARCH=amd64 go build  -o bin/mop-uploader-windows-amd64.exe .
	GOOS=linux GOARCH=amd64 go build  -o bin/mop-uploader-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build  -o bin/mop-uploader-darwin-amd64 .