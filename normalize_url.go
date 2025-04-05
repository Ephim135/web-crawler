package main

import (
	"errors"
	"net/url"
	"strings"
)

func normalizeURL(input string) (string, error) {
	parsed, err := url.Parse(input)
	if err != nil || parsed.Host == "" {
		return "", errors.New("invalid URL")
	}

	// Remove scheme
	host := parsed.Hostname()

	if strings.Contains(host, "..") {
		return "", errors.New("invalid hostname")
	}

	// Remove "www." prefix if present
	host = strings.TrimPrefix(host, "www.")

	// Reconstruct URL without scheme, port, query, or fragment
	normalized := host + parsed.EscapedPath()

	// Remove trailing slash if present
	normalized = strings.TrimSuffix(normalized, "/")

	return normalized, nil
}
