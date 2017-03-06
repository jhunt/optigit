package main

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/jhunt/go-db"
)

func findRepo(d db.DB, org, repo string) (int, error) {
	r, err := d.Query(`SELECT id FROM repos WHERE org = ? AND name = ?`, org, repo)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	if r.Next() {
		var id int
		err = r.Scan(&id)
		return id, err
	}
	return 0, nil
}

func importRepo(d db.DB, org, repo string) error {
	id, err := findRepo(d, org, repo)
	if err != nil || id > 0 {
		return err
	}

	fmt.Printf("import :: importing repo %s/%s into database\n", org, repo)
	err = d.Exec(`INSERT INTO repos (org, name, included) VALUES (?, ?, 1)`, org, repo)
	if err != nil {
		return err
	}
	_, err = findRepo(d, org, repo)
	return err
}

func clearIssues(d db.DB, org, repo string) error {
	id, err := findRepo(d, org, repo)
	if err != nil {
		return err
	}

	fmt.Printf("import :: deleting all cached issues for %s/%s from database\n", org, repo)
	return d.Exec(`DELETE FROM issues WHERE repo_id = ?`, id)
}

func importIssue(d db.DB, org, repo string, issue *github.Issue) error {
	id, err := findRepo(d, org, repo)
	if err != nil {
		return err
	}

	var created int64 = -1
	var updated int64 = -1

	if issue.CreatedAt != nil {
		created = issue.CreatedAt.Unix()
	}
	if issue.UpdatedAt != nil {
		updated = issue.UpdatedAt.Unix()
	}

	ppl := make([]string, len(issue.Assignees))
	for i, who := range issue.Assignees {
		ppl[i] = *who.Login
	}

	user := ""
	if issue.User != nil {
		user = *issue.User.Login
	}

	fmt.Printf("import :: importing issue '#%d - %s' for %s/%s into database\n", *issue.Number, *issue.Title, org, repo)
	return d.Exec(`INSERT INTO issues (id, repo_id, created_at, updated_at, reporter, assignees, title) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		*issue.Number, id, created, updated, user, strings.Join(ppl, ","), *issue.Title)
}

func clearPulls(d db.DB, org, repo string) error {
	id, err := findRepo(d, org, repo)
	if err != nil {
		return err
	}

	fmt.Printf("import :: deleting all cached pull requests for %s/%s from database\n", org, repo)
	return d.Exec(`DELETE FROM pulls WHERE repo_id = ?`, id)
}

func importPull(d db.DB, org, repo string, pull *github.PullRequest) error {
	id, err := findRepo(d, org, repo)
	if err != nil {
		return err
	}

	var created int64 = -1
	var updated int64 = -1

	if pull.CreatedAt != nil {
		created = pull.CreatedAt.Unix()
	}
	if pull.UpdatedAt != nil {
		updated = pull.UpdatedAt.Unix()
	}

	ppl := make([]string, len(pull.Assignees))
	for i, who := range pull.Assignees {
		ppl[i] = *who.Login
	}

	user := "EMPTY"
	if pull.User != nil {
		user = *pull.User.Login
	}

	fmt.Printf("import :: importing pull request '#%d - %s' for %s/%s into database\n", *pull.Number, *pull.Title, org, repo)
	return d.Exec(`INSERT INTO pulls (id, repo_id, created_at, updated_at, reporter, assignees, title) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		*pull.Number, id, created, updated, user, strings.Join(ppl, ","), *pull.Title)
}
