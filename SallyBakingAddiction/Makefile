BINARY_NAME=sally.exe

build:
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}-windows

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all