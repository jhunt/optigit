package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jhunt/go-db"
)

type RepoWatch struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func UpdateRepos(d db.DB, lst []RepoWatch) error {
	include := make([]interface{}, 0)

	for _, watch := range lst {
		id, err := strconv.ParseInt(watch.Name, 10, 0)
		if err != nil {
			return err
		}
		if watch.Value == "on" {
			include = append(include, int(id))
		}
	}

	if len(include) > 0 {
		err := d.Exec(`UPDATE repos SET included = 1 WHERE id IN (?`+strings.Repeat(`,?`, len(include)-1)+`)`, include...)
		if err != nil {
			return err
		}
		err = d.Exec(`UPDATE repos SET included = 0 WHERE id NOT IN (?`+strings.Repeat(`,?`, len(include)-1)+`)`, include...)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadRepos(d db.DB) ([]Repository, error) {
	repos, err := d.Query(`SELECT id, org, name, included FROM repos`)
	if err != nil {
		return nil, err
	}
	defer repos.Close()

	l := make([]Repository, 0)
	for repos.Next() {
		var (
			id, incl  int
			org, name string
		)
		err = repos.Scan(&id, &org, &name, &incl)
		if err != nil {
			return nil, err
		}

		l = append(l, Repository{
			ID:       id,
			Org:      org,
			Name:     name,
			Included: incl == 1,
		})
	}

	return l, nil
}

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
			ID:   id,
			Org:  org,
			Name: name,
		}
		issues, err := d.Query(`SELECT id, title, reporter, assignees, created_at, updated_at FROM issues WHERE repo_id = ?`, id)
		if err != nil {
			return nil, err
		}
		defer issues.Close()

		repo.Issues = make([]Issue, 0)
		for issues.Next() {
			var (
				number           int
				title, reporter, assignees string
				created, updated int
			)
			err = issues.Scan(&number, &title, &reporter, &assignees, &created, &updated)
			if err != nil {
				return nil, err
			}

			repo.Issues = append(repo.Issues, Issue{
				Number:  number,
				Title:   title,
				URL:     fmt.Sprintf("https://github.com/%s/%s/issues/%d", repo.Org, repo.Name, number),
				Created: created,
				Updated: updated,

				Reporter: reporter,
				Assignees: split(assignees),
			})
		}

		pulls, err := d.Query(`SELECT id, title, reporter, assignees, created_at, updated_at FROM pulls WHERE repo_id = ?`, id)
		if err != nil {
			return nil, err
		}
		defer pulls.Close()

		repo.PullRequests = make([]PullRequest, 0)
		for pulls.Next() {
			var (
				number           int
				title, reporter, assignees string
				created, updated int
			)
			err = pulls.Scan(&number, &title, &reporter, &assignees, &created, &updated)
			if err != nil {
				return nil, err
			}

			repo.PullRequests = append(repo.PullRequests, PullRequest{
				Number:  number,
				Title:   title,
				URL:     fmt.Sprintf("https://github.com/%s/%s/pull/%d", repo.Org, repo.Name, number),
				Created: created,
				Updated: updated,

				Reporter: reporter,
				Assignees: split(assignees),
			})
		}

		health[fmt.Sprintf("%s/%s", repo.Org, repo.Name)] = repo
	}

	return health, nil
}
