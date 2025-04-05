package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name           string
		inputURL       string
		expected       string
		expectedToFail bool
	}{
		{
			name:           "remove scheme",
			inputURL:       "https://blog.boot.dev/path",
			expected:       "blog.boot.dev/path",
			expectedToFail: false,
		},
		{
			name:           "remove end",
			inputURL:       "http://blog.boot.dev/path/",
			expected:       "blog.boot.dev/path",
			expectedToFail: false,
		},
		{
			name:           "remove end2",
			inputURL:       "https://blog.boot.dev/path/",
			expected:       "blog.boot.dev/path",
			expectedToFail: false,
		},
		{
			name:           "this test should fail",
			inputURL:       "http:///blog.boot.dev/path//",
			expected:       "error",
			expectedToFail: true,
		},
		{
			name:           "remove www prefix",
			inputURL:       "https://www.example.com/path",
			expected:       "example.com/path",
			expectedToFail: false,
		},
		{
			name:           "remove query parameters",
			inputURL:       "http://blog.boot.dev/path?query=123",
			expected:       "blog.boot.dev/path",
			expectedToFail: false,
		},
		{
			name:           "remove fragment/hash",
			inputURL:       "https://blog.boot.dev/path#section1",
			expected:       "blog.boot.dev/path",
			expectedToFail: false,
		},
		{
			name:           "remove port number",
			inputURL:       "http://blog.boot.dev:8080/path",
			expected:       "blog.boot.dev/path",
			expectedToFail: false,
		},
		{
			name:           "invalid URL format",
			inputURL:       "http://..blog.boot.dev/path",
			expected:       "error",
			expectedToFail: true,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expectedToFail {
				actual, err := normalizeURL(tc.inputURL)
				if err != nil {
					t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
					return
				}
				if actual != tc.expected {
					t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
				}
			} else {
				_, err := normalizeURL(tc.inputURL)
				if err == nil {
					t.Errorf("Test %v - %s FAIL: expected an error but got none", i, tc.name)
				}
				return // Exit early since we expect an error
			}
		})
	}
}
