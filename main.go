package main

import (
	"fmt"
	"os"
)

func main() {
	if os.Getenv("GITHUB_TOKEN") == "" {
		fmt.Printf("no GITHUB_TOKEN supplied via environment.\n")
		os.Exit(1)
	}
	if os.Getenv("ORGS") == "" {
		fmt.Printf("no ORGS supplied via environment.\n")
		os.Exit(1)
	}
	if os.Getenv("DATABASE") == "" && os.Getenv("VCAP_SERVICES") == "" {
		fmt.Printf("no database supplied via DATABASE or VCAP_SERVICES.\n")
		os.Exit(1)
	}

	bind := bindto()
	fmt.Printf("listening on %s\n", bind)
	if os.Getenv("DATABASE") != "" {
		fmt.Printf("using database %s\n", os.Getenv("DATABASE"))
	} else {
		_, dsn, _ := vcapdb(os.Getenv("VCAP_SERVICES"))
		fmt.Printf("using vcap %s database\n", dsn)
	}
	d, err := database()
	if err != nil {
		fmt.Printf("could not reach database: %v\n", err)
	}
	fmt.Printf("using github token %s\n", os.Getenv("GITHUB_TOKEN"))
	RunAPI(bind, d)
}
