GOOS=linux GOARCH=amd64 go build -o agent . && scp ./agent root@fedora-server:/home/maddox/agent
