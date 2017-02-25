package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/jhunt/optigit/static"
)

func RunAPI(bind string) {
	http.HandleFunc("/v1/scrape", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(404)
			fmt.Fprintf(w, "404 not found\n")
			return
		}
		fmt.Printf("Starting scrape...\n")

		d, err := database()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed connecting to backend database: %s\n", err)
		}

		orgs := strings.Split(os.Getenv("ORGS"), " ")
		err = Scrape(os.Getenv("GITHUB_TOKEN"), d, orgs...)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed: %s\n", err)
			return
		}
		w.WriteHeader(204)
		fmt.Printf("Finished scraping Github data!\n")
	})

	http.HandleFunc("/v1/health", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			w.WriteHeader(404)
			fmt.Fprintf(w, "404 not found\n")
			return
		}

		d, err := database()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed connecting to backend database: %s\n", err)
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

	http.HandleFunc("/v1/repos", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" && req.Method != "POST" {
			w.WriteHeader(404)
			fmt.Fprintf(w, "404 not found\n")
			return
		}

		d, err := database()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed connecting to backend database: %s\n", err)
			return
		}

		if req.Method == "GET" {
			repos, err := ReadRepos(d)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "failed connecting to backend database: %s\n", err)
				return
			}

			b, err := json.Marshal(repos)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "failed formatting JSON: %s\n", err)
				return
			}

			w.Header().Set("Content-type", "application/json")
			fmt.Fprintf(w, "%s\n", string(b))
			return
		}

		if req.Method == "POST" {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "failed reading request body: %s\n", err)
				return
			}

			var updates []RepoWatch
			err = json.Unmarshal(b, &updates)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "failed reading request JSON: %s\n", err)
				return
			}

			err = UpdateRepos(d, updates)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "failed applying updates: %s\n", err)
				return
			}

			w.WriteHeader(200)
			w.Header().Set("Content-type", "application/json")
			fmt.Fprintf(w, "{}")
			return
		}
	})

	http.Handle("/", static.Handler{})

	err := http.ListenAndServe(bind, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "http server exited: %s\n", err)
	}
}
