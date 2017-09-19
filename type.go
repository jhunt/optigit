package main

type PullRequest struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	Created   int      `json:"created"`
	Updated   int      `json:"updated"`
	Reporter  string   `json:"reporter"`
	Assignees []string `json:"assignees"`
}

type Issue struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	Created   int      `json:"created"`
	Updated   int      `json:"updated"`
	Reporter  string   `json:"reporter"`
	Assignees []string `json:"assignees"`
}

type Repository struct {
	ID           int           `json:"id"`
	Included     bool          `json:"included"`
	Org          string        `json:"org"`
	Name         string        `json:"name"`
	PullRequests []PullRequest `json:"pulls"`
	Issues       []Issue       `json:"issues"`
}

type Health struct {
	Repos  map[string]Repository `json:"repos"`
	Ignore string                `json:"ignore"`
}
