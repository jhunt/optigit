package main

import (
	"context"
	"strings"

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

func is404(err error) bool {
	return strings.Contains(err.Error(), "404 ")
}

func (g *Github) IssuesFor(who, repo string) ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var l []*github.Issue
	for {
		page, resp, err := g.Client.Issues.ListByRepo(g.Context, who, repo, opt)
		if err != nil {
			if is404(err) {
				return nil, nil
			}
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

func (g *Github) PullsFor(who, repo string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var l []*github.PullRequest
	for {
		page, resp, err := g.Client.PullRequests.List(g.Context, who, repo, opt)
		if err != nil {
			if is404(err) {
				return nil, nil
			}
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

func (g *Github) ScrapeRepos(d db.DB, who string) error {
	repos, err := g.ReposFor(who)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		err = importRepo(d, who, repo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Github) ScrapeIssues(d db.DB, who, repo string) error {
	issues, err := g.IssuesFor(who, repo)
	if err != nil {
		return err
	}

	err = clearIssues(d, who, repo)
	if err != nil {
		return err
	}
	for _, issue := range issues {
		err = importIssue(d, who, repo, issue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Github) ScrapePulls(d db.DB, who, repo string) error {
	pulls, err := g.PullsFor(who, repo)
	if err != nil {
		return err
	}

	err = clearPulls(d, who, repo)
	if err != nil {
		return err
	}
	for _, pull := range pulls {
		err = importPull(d, who, repo, pull)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Github) ReposFor(who string) ([]string, error) {
	oopt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	uopt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var l []*github.Repository
	for {
		page, resp, err := g.Client.Repositories.ListByOrg(g.Context, who, oopt)
		if err != nil {
			page, resp, err = g.Client.Repositories.List(g.Context, who, uopt)
			if err != nil {
				return nil, err
			}
		}
		l = append(l, page...)
		if resp.NextPage == 0 {
			break
		}
		oopt.ListOptions.Page = resp.NextPage
		uopt.ListOptions.Page = resp.NextPage
	}

	names := make([]string, len(l))
	for i := range l {
		names[i] = *l[i].Name
	}
	return names, nil
}

func Scrape(token string, d db.DB, users_or_orgs ...string) error {
	g := NewGithub(token)
	for _, who := range users_or_orgs {
		err := g.ScrapeRepos(d, who)
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
			who, repo string
		)
		err = r.Scan(&who, &repo)
		if err != nil {
			return err
		}
		repos = append(repos, [2]string{who, repo})
	}
	r.Close()

	for _, s := range repos {
		who := s[0]
		repo := s[1]

		err = g.ScrapeIssues(d, who, repo)
		if err != nil {
			return err
		}

		err = g.ScrapePulls(d, who, repo)
		if err != nil {
			return err
		}
	}

	DedupePullRequests(d)
	return nil
}
