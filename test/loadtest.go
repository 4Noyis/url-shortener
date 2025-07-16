package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Test configuration
	const (
		numRequests = 10000
		targetURL   = "http://localhost:8080/api/v1/q1m" // Replace with actual short URL
	)

	// Results tracking
	var (
		successCount int
		errorCount   int
		totalTime    time.Duration
		mutex        sync.Mutex
		wg           sync.WaitGroup
	)

	// Start timing
	startTime := time.Now()

	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Launch concurrent requests
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			reqStart := time.Now()
			resp, err := client.Get(targetURL)
			reqDuration := time.Since(reqStart)

			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				errorCount++
				fmt.Printf("Request %d failed: %v\n", requestID, err)
			} else {
				successCount++
				resp.Body.Close()
				fmt.Printf("Request %d: Status %d, Duration: %v\n", requestID, resp.StatusCode, reqDuration)
			}

			totalTime += reqDuration
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()

	// Calculate results
	overallDuration := time.Since(startTime)
	avgResponseTime := totalTime / time.Duration(numRequests)

	// Print results
	fmt.Printf("\n=== Load Test Results ===\n")
	fmt.Printf("Total requests: %d\n", numRequests)
	fmt.Printf("Successful requests: %d\n", successCount)
	fmt.Printf("Failed requests: %d\n", errorCount)
	fmt.Printf("Success rate: %.2f%%\n", float64(successCount)/float64(numRequests)*100)
	fmt.Printf("Overall test duration: %v\n", overallDuration)
	fmt.Printf("Average response time: %v\n", avgResponseTime)
	fmt.Printf("Requests per second: %.2f\n", float64(numRequests)/overallDuration.Seconds())
}
