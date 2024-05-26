BINARY_NAME=img_proc
BINARY_LINUX=${BINARY_NAME}-linux
BINARY_WINDOWS=${BINARY_NAME}-windows

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_LINUX} ./cmd/img_prog.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_WINDOWS} ./cmd/img_prog.go

run: build
	chmod +x ${BINARY_LINUX}
	./${BINARY_LINUX}

clean:
	go clean
	rm -f ${BINARY_LINUX} ${BINARY_WINDOWS}
