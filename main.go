package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)


func main() {
	now := time.Now()
	args := os.Args[1:]
	lengthArgs := len(args)
	if lengthArgs < 3 {
		fmt.Println("no website provided")
		os.Exit(1)	
	}
	if lengthArgs > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", args[0])

	URL := args[0]
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("error second argument should be int")
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("error third argument should be int")
	}

	cfg := newConfig(URL, maxConcurrency, maxPages)

	// the wait Group counts how many go routines there are
	cfg.wg.Add(1)
	go cfg.crawlPage(args[0])
	// this line is like a blocking wall for the code untill the wg is 0
	// meaning all go routines called Done
	cfg.wg.Wait()

	printReport(cfg.pages, URL) 
	fmt.Println(time.Since(now))
}

func printReport(pages map[string]int, baseURL string) {
	type kv struct {
		Key string
		Value int
	}

	var sortedData []kv
	for k, v := range pages {
		sortedData = append(sortedData, kv{k, v})
	}

	sort.Slice(sortedData, func(i, j int) bool {
		if sortedData[i].Value != sortedData[j].Value {
			return sortedData[i].Value > sortedData[j].Value
		}
		return sortedData[i].Key < sortedData[j].Key
	})

	fmt.Println("============================")
	fmt.Printf("  REPORT for %s\n", baseURL)
	fmt.Println("============================")
	for _, item := range sortedData {
		fmt.Printf("Found %d internal links to %s\n", item.Value, item.Key)
	}
}