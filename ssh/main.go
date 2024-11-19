package main

import (
	"os"
	"time"
)

func main() {
	client, err := OpenClient(SessionOpts{
		Username:      "root",
		Password:      "subluxation",
		ServerAddress: "fedora-server",
	})

	if err != nil {
		panic(err)
	}

	file, err := os.Create("test.txt")

	if err != nil {
		panic(err)
	}

	err = client.RunWithOutputStream("ls -la", file, file)

	if err != nil {
		panic(err)
	}

	for {
		err = client.RunWithOutputStream("ls -la", os.Stdout, os.Stderr)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}

	client.Disconnect()
}
