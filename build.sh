#!/bin/bash
source .env

echo $SERVER

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build -o neo-moonchan .

# should need further test.
ssh $SERVER "killall neo-moonchan"
scp neo-moonchan $SERVER:~/
ssh $SERVER "nohup ~/neo-moonchan &"
# scp neo-moonchan $SERVER:~/neo-moonchan/
# scp .env $SERVER:~/neo-moonchan/

rm neo-moonchan

