# go-chat

A go based chat app with the simple C-S achitecture from scratch.

## Usage
go build -o server src/main.go src/server.go src/user.go

go build -o client src/client.go

Open a terminal:
./server

Open seperate terminals for clients:
./client -ip 127.0.0.1 -port 8080