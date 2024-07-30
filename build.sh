#!/bin/bash
source .env

echo $SERVER

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build -o neo-moonchan .

scp neo-moonchan $SERVER:~/
# scp neo-moonchan $SERVER:~/neo-moonchan/
# scp .env $SERVER:~/neo-moonchan/

rm neo-moonchan

