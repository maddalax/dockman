htmgo setup && GOOS=linux GOARCH=amd64 go build --tags prod -o dockside . && scp ./dockside root@fedora-server:/home/maddox/dockside
cd ./cmd/agent && GOOS=linux GOARCH=amd64 go build -o agent . && scp ./agent root@fedora-server:/home/maddox/agent
