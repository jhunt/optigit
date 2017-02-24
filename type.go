package main

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
	ID           int           `json:"id"`
	Included     bool          `json:"included"`
	Org          string        `json:"org"`
	Name         string        `json:"name"`
	Branches     []string      `json:"branches"`
	PullRequests []PullRequest `json:"pulls"`
	Issues       []Issue       `json:"issues"`
}

type Health map[string]Repository
