package main

import (
	"fmt"
	"reflect"
	"testing"
)
func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
			<html>
				<body>
					<a href="/path/one">
						<span>Boot.dev</span>
					</a>
					<a href="https://other.com/path/one">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "no anchor tags",
			inputURL: "https://example.com",
			inputBody: `
			<html>
				<body>
					<p>No links here</p>
				</body>
			</html>
			`,
			expected: []string{},
		},
		{
			name:     "no href attribute",
			inputURL: "https://example.com",
			inputBody: `
			<html>
				<body>
					<a>No href here</a>
				</body>
			</html>
			`,
			expected: []string{},
		},
		{
			name:     "malformed URL",
			inputURL: "https://example.com",
			inputBody: `
			<html>
				<body>
					<a href="htp://invalid-url">Invalid</a>
				</body>
			</html>
			`,
			expected: []string{}, // the malformed URL should be ignored
		},
		{
			name:     "empty HTML",
			inputURL: "https://example.com",
			inputBody: "",
			expected: []string{}, // no links should be found
		},
		{
			name:     "empty href",
			inputURL: "https://example.com",
			inputBody: `
			<html>
				<body>
					<a href="">Empty Link</a>
				</body>
			</html>
			`,
			expected: []string{"https://example.com"}, // empty href should resolve to the base URL
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
			}
			// reflect.DeepEqual did not compare correctly because the empty slice have been different
			if !reflect.DeepEqual(got, append([]string(nil), tc.expected...)) {
				t.Errorf("Test %v - %s FAIL: expected: %v, actual: %v", i, tc.name, tc.expected, got)
			} else {
				fmt.Printf("Test %v success\n", i)
			}
		})
	}
}