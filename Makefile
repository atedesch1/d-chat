all: client

client: cmd/client/main.go
	go build -o bin/client cmd/client/*


clean:
	rm -rf bin/*