GOOS=linux GOARCH=amd64 go build -o keyring . && scp ./keyring root@fedora-server:/home/maddox/keyring
