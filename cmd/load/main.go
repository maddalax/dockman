package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	hostname       = "https://fedora.htmgo.dev"
	requestsPerSec = 800
	duration       = 1 * time.Minute // Run the test for 10 seconds
)

var urls = []string{
	"/docs",
	"/examples",
	"/",
	"/html-to-go",
}

func main() {
	var wg sync.WaitGroup
	var successCount, failureCount int64

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	startTime := time.Now()
	endTime := startTime.Add(duration)

	fmt.Printf("Starting load test at %v\n", startTime)

	for time.Now().Before(endTime) {
		for i := 0; i < requestsPerSec/len(urls); i++ {
			for _, url := range urls {
				wg.Add(1)
				go func(url string) {
					defer wg.Done()
					resp, err := client.Get(hostname + url)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						failureCount++
						return
					}
					resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						successCount++
					} else {
						fmt.Printf("Non-200 response: %v\n", resp.Status)
						failureCount++
					}
				}(url)
			}
		}
		time.Sleep(time.Second)
	}

	wg.Wait()

	endTime = time.Now()
	fmt.Printf("Load test completed at %v\n", endTime)
	fmt.Printf("Total requests: %d\n", successCount+failureCount)
	fmt.Printf("Successful requests: %d\n", successCount)
	fmt.Printf("Failed requests: %d\n", failureCount)
}
