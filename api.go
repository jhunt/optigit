package main

import (
	"fmt"

	"github.com/jhunt/go-db"
)

func ReadInformation(d db.DB) (Health, error) {
	health := make(map[string]Repository)

	repos, err := d.Query(`SELECT id, org, name FROM repos WHERE included = 1`)
	if err != nil {
		return nil, err
	}
	defer repos.Close()

	for repos.Next() {
		var (
			id        int
			org, name string
		)
		err = repos.Scan(&id, &org, &name)
		if err != nil {
			return nil, err
		}

		repo := Repository{
			Org:  org,
			Name: name,
		}
		issues, err := d.Query(`SELECT id, title, assignees, created_at, updated_at FROM issues WHERE repo_id = ?`, id)
		if err != nil {
			return nil, err
		}
		defer issues.Close()

		repo.Issues = make([]Issue, 0)
		for issues.Next() {
			var (
				number           int
				title, assignees string
				created, updated int
			)
			err = issues.Scan(&number, &title, &assignees, &created, &updated)
			if err != nil {
				return nil, err
			}

			repo.Issues = append(repo.Issues, Issue{
				Number:  number,
				Title:   title,
				URL:     fmt.Sprintf("https://github.com/%s/%s/issues/%d", repo.Org, repo.Name, number),
				Created: created,
				Updated: updated,

				Assignees: split(assignees),
			})
		}

		pulls, err := d.Query(`SELECT id, title, assignees, created_at, updated_at FROM pulls WHERE repo_id = ?`, id)
		if err != nil {
			return nil, err
		}
		defer pulls.Close()

		repo.PullRequests = make([]PullRequest, 0)
		for pulls.Next() {
			var (
				number           int
				title, assignees string
				created, updated int
			)
			err = pulls.Scan(&number, &title, &assignees, &created, &updated)
			if err != nil {
				return nil, err
			}

			repo.PullRequests = append(repo.PullRequests, PullRequest{
				Number:  number,
				Title:   title,
				URL:     fmt.Sprintf("https://github.com/%s/%s/pull/%d", repo.Org, repo.Name, number),
				Created: created,
				Updated: updated,

				Assignees: split(assignees),
			})
		}

		health[fmt.Sprintf("%s/%s", repo.Org, repo.Name)] = repo
	}

	return health, nil
}

