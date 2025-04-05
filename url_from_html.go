package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	var URLs []string
	var rawHrefValues []string

	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return URLs, fmt.Errorf("error parsing htmlBody: %v", err) 
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return URLs, fmt.Errorf("error parsing rawBaseURL: %v", err)		
	}

	// loop all Descendants and get the values of the href values
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					rawHrefValues = append(rawHrefValues, a.Val)
					break		
				}
			}
		}
	}

	// make relative Paths(/some/path) into urls with the rawBaseUrl and
	// add them to URLs normal URLs too and malformed ones get sroted out
	for _,rawURL := range rawHrefValues {
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			continue
		}

		
		if parsedURL.Scheme == "" {
			parsedURL = baseURL.ResolveReference(parsedURL)
		}
		
		// check if schema is http or https
		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			continue
		}

		URLs = append(URLs, parsedURL.String())
	}
	// make Urls from href values
	// find out if relative Path or url
	return URLs, nil
}
