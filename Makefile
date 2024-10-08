gen:
	protoc --proto_path=proto proto/*.proto --go_out=./pkg/pb --go_opt=paths=source_relative --go-grpc_out=./pkg/pb --go-grpc_opt=paths=source_relative

clean:
	rm ./pkg/pb/*.go

evans:
	evans -r repl -p 8080

servermongo:
	go run cmd/server/main.go --port 8081 --dsn="mongodb://username:password@127.0.0.1:27018/pcbook?parseTime=true&authSource=pcbook&authMechanism=SCRAM-SHA-256" --db="pcbook"

server1:
	go run cmd/server/main.go --port 50001

server2:
	go run cmd/server/main.go --port 50002

server:
	 go run cmd/server/main.go --port 8080

servertls:
	 go run cmd/server/main.go --port 4433 --tls=true

client:
	go run cmd/client/main.go --address 0.0.0.0:8080

clienttls:
	go run cmd/client/main.go --address 0.0.0.0:4433 --tls=true

test:
	go test -v -coverprofile cover.out ./...

cert:
	cd ./cert && sh gen.sh

.PHONY: gen clean evans server client test cert servertls clienttls server1 server2
