all: client

client: cmd/client/main.go
	go build -o bin/client cmd/client/*

zookeeper:
	docker-compose up

clean:
	rm -rf bin/*