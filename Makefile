all: app protogen client

app: cmd/app/main.go
	go build -o bin/app cmd/app/*

client: cmd/client/main.go
	go build -o bin/client cmd/client/*

zookeeper:
	docker-compose up

protogen:
	protoc --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--proto_path=proto proto/*.proto

clean:
	rm -rf bin/*
	rm pb/*.go