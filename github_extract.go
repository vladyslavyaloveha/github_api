package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v35/github"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

func getEnvVar(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid key or type assertion: key=%s, value=%#v", key, value)
	}
	return value
}

func getGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func getRepository(client *github.Client, query string) (*Repository, error) {
	opts := &github.SearchOptions{Sort: "stars", Order: "desc",
		ListOptions: github.ListOptions{Page: 1, PerPage: 1}}
	repos, response, err := client.Search.Repositories(context.Background(), query, opts)
	if err != nil {
		return nil, fmt.Errorf("get repository %s: ", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get repository status %d: ", response.StatusCode)
	}
	if len(repos.Repositories) == 0 {
		return nil, fmt.Errorf("repository not found: query: %s", query)
	}

	repository := &Repository{
		ID:          repos.Repositories[0].ID,
		Owner:       repos.Repositories[0].Owner,
		Name:        repos.Repositories[0].Name,
		FullName:    repos.Repositories[0].FullName,
		Description: repos.Repositories[0].Description,
		CreatedAt:   repos.Repositories[0].CreatedAt,
		Language:    repos.Repositories[0].Language,
		GitURL:      repos.Repositories[0].GitURL,
	}
	return repository, nil
}

func getIssues(client *github.Client, query string, filter IssueFilter) ([]Issue, error) {
	repository, err := getRepository(client, query)
	if err != nil {
		return nil, fmt.Errorf("get issues: %s", err)
	}
	options := &github.IssueListByRepoOptions{State: filter.State}
	issues_, response, err := client.Issues.ListByRepo(context.Background(), *repository.Owner.Login,
		*repository.Name, options)
	if err != nil {
		return nil, fmt.Errorf("get issues: %s", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get issues status %d: ", response.StatusCode)
	}

	var issues []Issue
	for _, issue := range issues_ {
		issues = append(issues, Issue{
			ID:            issue.ID,
			State:         issue.State,
			Title:         issue.Title,
			Body:          issue.Body,
			User:          issue.User,
			Assignee:      issue.Assignee,
			ClosedAt:      issue.ClosedAt,
			CreatedAt:     issue.CreatedAt,
			UpdatedAt:     issue.UpdatedAt,
			ClosedBy:      issue.ClosedBy,
			URL:           issue.URL,
			RepositoryURL: issue.RepositoryURL,
			Assignees:     issue.Assignees,
		})
	}
	return issues, nil
}

func getIssueState(ctx iris.Context) string {
	state := ctx.URLParamDefault("state", "all")
	switch state {
	case "all":
		return state
	case "open":
		return state
	case "closed":
		return state
	default:
		return "all"
	}
}

func getCommits(client *github.Client, query string, filter CommitFilter) ([]Commit, error) {
	repository, err := getRepository(client, query)
	if err != nil {
		return nil, fmt.Errorf("get issues: %s", err)
	}

	options := &github.CommitsListOptions{Author: filter.Author, Since: filter.Since, Until: filter.Until}
	commits_, response, err := client.Repositories.ListCommits(context.Background(), *repository.Owner.Login,
		*repository.Name, options)
	if err != nil {
		return nil, fmt.Errorf("get commits: %s", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get commits status %d: ", response.StatusCode)
	}

	var commits []Commit
	for _, commit := range commits_ {
		commits = append(commits, Commit{
			SHA:       commit.SHA,
			Author:    commit.GetCommit().Author,
			Committer: commit.GetCommit().Committer,
			Message:   commit.GetCommit().Message,
			URL:       commit.GetCommit().URL,
		})
	}
	return commits, nil
}

func getCommitAuthor(ctx iris.Context) string {
	author := ctx.URLParamDefault("author", "")
	return author
}
