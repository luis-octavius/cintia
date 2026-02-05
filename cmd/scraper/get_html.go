package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var Client = &http.Client{}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating the request: %v", err)
	}

	req.Header.Set("Content-Type", "text/html")
	req.Header.Set("Crawler", "OctaviusCrawler 1.0")

	res, err := Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error doing the request: %v", err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading from body: %v", err)
	}

	return string(body), nil
}
