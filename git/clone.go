package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
)

func main() {
	_, err := git.PlainClone("/tmp/foo2", false, &git.CloneOptions{
		URL:      "https://github.com/maddalax/paas",
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: "paas",
			Password: "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",
		},
	})
	if err != nil {
		panic(err)
	}
}
