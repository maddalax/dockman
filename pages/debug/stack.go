package debug

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
	"regexp"
	"runtime"
	"strings"
)

func PrintStack(ctx *h.RequestContext) *h.Page {
	buf := make([]byte, 1<<20) // 1 MB buffer size
	n := runtime.Stack(buf, true)
	data := fmt.Sprintf("%s\n", buf[:n])
	lines := strings.Split(data, "\n")

	// Create a regular expression to identify lines that start with 'goroutine'
	goroutineRegex := regexp.MustCompile(`^goroutine \d+`)

	var blocks [][]string
	var currentBlock []string

	var shouldIncludeBlock = func(block []string) bool {
		for _, s := range block {
			if strings.Contains(s, "github.com/nats-io/nats-server/v2/server") {
				return false
			}
		}

		return true
	}

	// Iterate over each line and group into blocks
	for _, line := range lines {
		if goroutineRegex.MatchString(line) {
			if len(currentBlock) > 0 && shouldIncludeBlock(currentBlock) {
				blocks = append(blocks, currentBlock)
			}
			currentBlock = []string{line}
		} else {
			currentBlock = append(currentBlock, line)
		}
	}
	if len(currentBlock) > 0 && shouldIncludeBlock(currentBlock) {
		blocks = append(blocks, currentBlock)
	}

	for i, block := range blocks {
		fmt.Println("Block", i)
		for _, line := range block {
			fmt.Println(line)
		}
		fmt.Println()
	}

	return h.EmptyPage()
}
