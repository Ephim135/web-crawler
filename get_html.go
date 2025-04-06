package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Gethtml(rawURL string) (string, error){
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("error get response from URL: %v", err) 
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return "", errors.New("400 status code from response")
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("content-type is not html: %s", contentType)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed read response Body: %v", err)
	}

	return string(data), nil
}