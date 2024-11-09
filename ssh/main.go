package main

import "os"

func main() {
	session, err := OpenSession(SessionOpts{
		Username:      "root",
		Password:      "subluxation",
		ServerAddress: "100.69.2.76",
	})

	if err != nil {
		panic(err)
	}

	file, err := os.Create("test.txt")

	if err != nil {
		panic(err)
	}

	err = session.RunWithOutputStream("ls -la", file, file)

	if err != nil {
		panic(err)
	}

	session.Disconnect()
}
