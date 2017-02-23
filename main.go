package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jhunt/go-db"
	_ "github.com/mattn/go-sqlite3"
)

type PullRequest struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	Created   int      `json:"created"`
	Updated   int      `json:"updated"`
	Assignees []string `json:"assignees"`
}

type Issue struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	Created   int      `json:"created"`
	Updated   int      `json:"updated"`
	Assignees []string `json:"assignees"`
}

type Repository struct {
	Org          string        `json:"org"`
	Name         string        `json:"name"`
	Branches     []string      `json:"branches"`
	PullRequests []PullRequest `json:"pulls"`
	Issues       []Issue       `json:"issues"`
}

type Health map[string]Repository

func split(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func getDB() db.DB {
	return db.DB{
		Driver: "sqlite3",
		DSN:    "./db/sq.db",
	}
}

func runServer() {
	http.HandleFunc("/v1/", func(w http.ResponseWriter, req *http.Request) {
		d := getDB()
		err := d.Connect()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed connecting to backend database: %s\n", err)
			return
		}
		if !d.Connected() {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed connecting to backend database: not connected\n")
			return
		}

		health, err := ReadInformation(d)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed connecting to backend database: %s\n", err)
			return
		}

		b, err := json.Marshal(health)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed formatting JSON: %s\n", err)
			return
		}

		w.Header().Set("Content-type", "application/json")
		fmt.Fprintf(w, "%s\n", string(b))
	})
	http.ListenAndServe(":3004", nil)
}

func runImport() {
	d := getDB()
	err := d.Connect()
	if err != nil {
		fmt.Printf("failed connecting to backend database: %s\n", err)
		return
	}
	if !d.Connected() {
		fmt.Printf("failed connecting to backend database: not connected\n")
		return
	}

	err = Scrape(os.Getenv("GITHUB_TOKEN"), d, "bolo")
	if err != nil {
		fmt.Printf("failed: %s\n", err)
	}
	return
}

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "web", "api":
			runServer()
			os.Exit(0)

		case "import":
			runImport()
			os.Exit(0)
		}
	}

	fmt.Printf("USAGE: optigit (api|import)\n")
	os.Exit(1)
}
