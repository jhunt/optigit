package main

import (
	"fmt"
	"os"
)

func main() {
	_, err := database()
	if err != nil {
		fmt.Printf("failed to talk to database on startup: %s\n", err)
		os.Exit(1)
	}
	if os.Getenv("GITHUB_TOKEN") == "" {
		fmt.Printf("no GITHUB_TOKEN supplied via environment.\n")
		os.Exit(1)
	}

	bind := bindto()
	fmt.Printf("listening on %s\n", bind)
	fmt.Printf("using database %s\n", os.Getenv("DATABASE"))
	fmt.Printf("using github token %s\n", os.Getenv("GITHUB_TOKEN"))
	RunAPI(bind)
}
