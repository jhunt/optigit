package main

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/jhunt/go-db"
	"golang.org/x/oauth2"
)

type Github struct {
	Client  *github.Client
	Context context.Context
}

func NewGithub(token string) *Github {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Github{
		Client:  github.NewClient(tc),
		Context: ctx,
	}
}

func (g *Github) IssuesFor(org, repo string) ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var l []*github.Issue
	for {
		page, resp, err := g.Client.Issues.ListByRepo(g.Context, org, repo, opt)
		if err != nil {
			return nil, err
		}
		l = append(l, page...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return l, nil
}

func (g *Github) PullsFor(org, repo string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var l []*github.PullRequest
	for {
		page, resp, err := g.Client.PullRequests.List(g.Context, org, repo, opt)
		if err != nil {
			return nil, err
		}
		l = append(l, page...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return l, nil
}

func (g *Github) ScrapeRepos(d db.DB, org string) error {
	repos, err := g.ReposFor(org)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		err = importRepo(d, org, repo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Github) ScrapeIssues(d db.DB, org, repo string) error {
	issues, err := g.IssuesFor(org, repo)
	if err != nil {
		return err
	}

	err = clearIssues(d, org, repo)
	if err != nil {
		return err
	}
	for _, issue := range issues {
		err = importIssue(d, org, repo, issue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Github) ScrapePulls(d db.DB, org, repo string) error {
	pulls, err := g.PullsFor(org, repo)
	if err != nil {
		return err
	}

	err = clearPulls(d, org, repo)
	if err != nil {
		return err
	}
	for _, pull := range pulls {
		err = importPull(d, org, repo, pull)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Github) ReposFor(org string) ([]string, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var l []*github.Repository
	for {
		page, resp, err := g.Client.Repositories.ListByOrg(g.Context, org, opt)
		if err != nil {
			return nil, err
		}
		l = append(l, page...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	names := make([]string, len(l))
	for i := range l {
		names[i] = *l[i].Name
	}
	return names, nil
}

func Scrape(token string, d db.DB, orgs ...string) error {
	g := NewGithub(token)
	for _, org := range orgs {
		err := g.ScrapeRepos(d, org)
		if err != nil {
			return err
		}
	}

	repos := make([][2]string, 0)
	r, err := d.Query(`SELECT org, name FROM repos WHERE included = 1`)
	if err != nil {
		return err
	}
	defer r.Close()

	for r.Next() {
		var (
			org, repo string
		)
		err = r.Scan(&org, &repo)
		if err != nil {
			return err
		}
		repos = append(repos, [2]string{org, repo})
	}
	r.Close()

	for _, s := range repos {
		org := s[0]
		repo := s[1]

		err = g.ScrapeIssues(d, org, repo)
		if err != nil {
			return err
		}

		err = g.ScrapePulls(d, org, repo)
		if err != nil {
			return err
		}
	}

	return nil
}
