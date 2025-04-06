package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages map[string]int 
	maxPages int
	baseURL *url.URL
	mu *sync.Mutex
	concurrencyControl chan struct{}
	wg *sync.WaitGroup
}

func newConfig(baseURL string, maxConcurrency int, maxPages int) *config {
	URL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &config{
		pages: make(map[string]int),
		maxPages: maxPages,
		baseURL: URL,
		mu: &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg: &sync.WaitGroup{},
	}
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	// takes a slot in the channel
	// Limits the amount of go routines that can run at a time
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		// frees the slot so the next go routine can start
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()
		
	// check if maximum amount of pages is reached if yes return
	if cfg.checkMaxPages() {
		return
	}

	currentUrl, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	// skip Other Websites
	if currentUrl.Hostname() != cfg.baseURL.Hostname() {
		return
	}
	fmt.Printf("crawling: %s\n", rawCurrentURL)

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error normalizedURL: %v", err)
	}

	// skip Website if it is not Visited for the first time
	if !cfg.addPageVisit(normalizedURL) {
		return
	}

	html, err := Gethtml(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	URLs, err := getURLsFromHTML(html, cfg.baseURL.String())
	if err != nil {
		fmt.Printf("Error getURLsFromHTTP: %v", err)
		return
	}
	fmt.Printf("found %d URLs\n", len(URLs))
	fmt.Println("-------------------------------------------------")

	for _, URL := range URLs {
		// somehow make this add to the wait group 
		cfg.wg.Add(1)
		go cfg.crawlPage(URL)

	}
}

func (cfg *config) addPageVisit(normalizedURL string) bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if cfg.pages[normalizedURL] > 0 {
		cfg.pages[normalizedURL]++
		fmt.Printf("Page already crawled %d Times\n", cfg.pages[normalizedURL])
		return false 
	}

	cfg.pages[normalizedURL] = 1
	fmt.Println("Page gets crawled first time")
	return true
}

func (cfg *config) checkMaxPages() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages) > cfg.maxPages
}