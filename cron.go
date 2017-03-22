package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jhunt/go-db"
)

func Cron(interval string, d db.DB) {
	re := regexp.MustCompile(`^(\d+(\.\d+)?)([hmd])$`)

	if interval == "" {
		interval = "1d"
	}

	when := 24 * time.Hour
	m := re.FindStringSubmatch(interval)
	if m == nil {
		fmt.Fprintf(os.Stderr, "'%s' is an invalid time spec (must be ##.#d or ##.#h)\n", interval)
		fmt.Fprintf(os.Stderr, "falling back to daily background job execution\n")
	} else {
		i, _ := strconv.Atoi(m[1])
		switch m[3] {
		case "m":
			fmt.Printf("setting background process to %d-minute frequency\n", i)
			when = time.Duration(i) * time.Minute
		case "h":
			fmt.Printf("setting background process to %d-hour frequency\n", i)
			when = time.Duration(i) * time.Hour
		case "d":
			fmt.Printf("setting background process to %d-day frequency\n", i)
			when = time.Duration(i) * 24 * time.Hour
		default:
			fmt.Fprintf(os.Stderr, "failed to detecet time units for frequency of background processes.\n")
			fmt.Fprintf(os.Stderr, "falling back to 1-day frequency, but this is not ideal...\n")
		}
	}

	waitforit := time.NewTicker(when).C
	for {
		<-waitforit
		fmt.Printf("performing regularly-scheduled background scrape")
		orgs := strings.Split(os.Getenv("ORGS"), " ")
		err := Scrape(os.Getenv("GITHUB_TOKEN"), d, orgs...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scrape failed: %s\n", err)
		} else {
			fmt.Printf("scrape finished successfully!\n")
		}
	}
}
