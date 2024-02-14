BINARY_NAME=toru
 
build:
	go mod tidy
	go build -o ${BINARY_NAME} ./cmd/toru
    

clean:
	go clean
	rm ${BINARY_NAME}
