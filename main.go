package main

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Data struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method,omitempty"`
	Body    string            `yaml:"body,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

// extractDomain extracts the domain from a given URL.
// It parses the provided full URL and returns the domain name without the port number, if present.
// For example:
//   - "https://www.example.com:8080/path" → "www.example.com"
//   - "http://sub.domain.com/page" → "sub.domain.com"
//
// Parameters:
// - fullURL: A string containing the full URL.
// Returns:
// - string: The extracted domain name.
// - error: An error if the URL is invalid or cannot be parsed.
func extractDomain(fullURL string) (string, error) {
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		//fmt.Printf("error extracting domain: %s\n", err)
		return "", err
	}
	// extract domain and remove port if present
	host := strings.Split(parsedURL.Host, ":")[0]
	return host, nil
}

// getStatusCode sends an HTTP request and returns the status code.
// It supports GET and POST requests, applying headers and sending a body if provided.
//
// Parameters:
// - data: A struct containing URL, method, optional headers, and optional body.
// Returns:
// - int: The HTTP status code.
// - error: Any error encountered.
func getStatusCode(data Data) (int, error) {
	client := &http.Client{Timeout: 5 * time.Second} // Set timeout

	var reqBody *bytes.Reader
	if data.Method == "POST" && data.Body != "" {
		reqBody = bytes.NewReader([]byte(data.Body))
	} else {
		reqBody = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(data.Method, data.URL, reqBody)
	if err != nil {
		//fmt.Printf("error creating request: %s\n", err)
		return 0, err
	}

	if data.Headers != nil {
		for key, value := range data.Headers {
			req.Header.Set(key, value)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		//fmt.Printf("error sending request: %s\n", err)
		return 0, err
	}
	defer res.Body.Close()

	return res.StatusCode, nil
}

func main() {

	// In-memory map to store domain statistics for program duration
	domainMap := make(map[string]map[string]int)

	yamlFile, err := os.ReadFile("sample2.yaml")
	if err != nil {
		fmt.Println("Error reading yaml file:", err)
		os.Exit(1)
	}

	var data []Data

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		fmt.Println("Error parsing yaml file:", err)
		os.Exit(1)
	}

	// Set default method to GET if not provided
	for i := range data {
		if data[i].Method == "" {
			data[i].Method = "GET"
		}
	}

	for {

		for _, req := range data {

			// Available components of req:
			// req.Name, req.URL, req.Method, req.Headers, req.Body

			fullDomain, _ := extractDomain(req.URL)

			if _, exists := domainMap[fullDomain]; !exists {
				domainMap[fullDomain] = map[string]int{
					"count": 1,
					"up":    0,
				}
			} else {
				domainMap[fullDomain]["count"] += 1
			}

			start := time.Now()

			statusCode, _ := getStatusCode(req)

			latency := time.Since(start).Milliseconds()

			if (statusCode >= 200 && statusCode <= 299) && latency < 500 {
				domainMap[fullDomain]["up"] += 1
			}

		}

		for domain, stats := range domainMap {

			upCount := stats["up"]
			totalCount := stats["count"]
			AvailabilityPercentage := math.Round((float64(upCount) / float64(totalCount)) * 100)

			fmt.Printf("%s has %d%% availability percentage\n", domain, int(AvailabilityPercentage))
		}

		time.Sleep(15 * time.Second)

	}
}
