package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// LoadTestRoundTripper is a custom RoundTripper optimized for load testing
type LoadTestRoundTripper struct {
	Transport *http.Transport
}

// NewLoadTestRoundTripper creates and configures a RoundTripper for load testing
func NewLoadTestRoundTripper() *LoadTestRoundTripper {
	transport := &http.Transport{

		// Optimize for load testing scenarios
		MaxIdleConns:        10000,
		MaxIdleConnsPerHost: 10000,
		MaxConnsPerHost:     10000,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &LoadTestRoundTripper{Transport: transport}
}

// RoundTrip delegates the request to the underlying Transport
func (r *LoadTestRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return r.Transport.RoundTrip(req)
}

const (
	hostname = "https://fedora.htmgo.dev"
	agents   = 50
	duration = 1 * time.Minute // Run the test for 10 seconds
)

var urls = []string{
	"/docs",
	"/examples",
	"/",
	"/html-to-go",
}

func main() {
	var wg sync.WaitGroup
	var successCount, failureCount, totalBytes atomic.Int64

	client := &http.Client{
		Transport: NewLoadTestRoundTripper(),
		Timeout:   10 * time.Second,
	}

	startTime := time.Now()
	endTime := startTime.Add(duration)

	fmt.Printf("Starting load test at %v\n", startTime)

	go func() {
		for {
			totalRequests := successCount.Load() + failureCount.Load()
			runningTime := time.Now().Sub(startTime)
			secondsElapsed := int64(runningTime.Seconds())
			if secondsElapsed == 0 {
				secondsElapsed = 1
			}
			perSecond := totalRequests / secondsElapsed
			timeLeft := endTime.Sub(time.Now()).Seconds()
			fmt.Printf("Successful: %d, Failed: %d, Per Second: %d, Seconds Left: %d, Bytes Read: %d \n", successCount.Load(), failureCount.Load(), perSecond, int(timeLeft), totalBytes.Load())
			time.Sleep(time.Second)
		}
	}()

	for i := 0; i < agents; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if time.Now().After(endTime) {
					return
				}
				time.Sleep(time.Millisecond * 10)
				url2 := urls[rand.Intn(len(urls))]
				url3, _ := url.Parse(hostname + url2)
				req := &http.Request{
					Method: http.MethodGet,
					URL:    url3,
					Header: make(http.Header),
				}
				req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyClient/1.0)")
				req.Header.Set("Accept", "*/*")
				req.Header.Set("Connection", "keep-alive")
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					failureCount.Add(1)
					return
				}
				readBody, _ := io.ReadAll(resp.Body)
				totalBytes.Add(int64(len(readBody)))
				if resp.StatusCode == http.StatusOK {
					successCount.Add(1)
				} else {
					fmt.Printf("Non-200 response: %v\n", resp.Status)
					failureCount.Add(1)
				}
				resp.Body.Close()
			}
		}()
	}

	wg.Wait()

	endTime = time.Now()
	fmt.Printf("Load test completed at %v\n", endTime)
	fmt.Printf("Total requests: %d\n", successCount.Load()+failureCount.Load())
	fmt.Printf("Successful requests: %d\n", successCount.Load())
	fmt.Printf("Failed requests: %d\n", failureCount.Load())
	fmt.Printf("Total bytes read: %d\n", totalBytes.Load())
}
