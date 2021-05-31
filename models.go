package main

import (
	"github.com/google/go-github/v35/github"
	"time"
)

type Repository struct {
	ID          *int64            `json:"id"`
	Owner       *github.User      `json:"owner"`
	Name        *string           `json:"name"`
	FullName    *string           `json:"full_name"`
	Description *string           `json:"description"`
	CreatedAt   *github.Timestamp `json:"created_at"`
	Language    *string           `json:"language"`
	GitURL      *string           `json:"git_url"`
}

type Issue struct {
	ID            *int64         `json:"id,omitempty"`
	State         *string        `json:"state,omitempty"`
	Title         *string        `json:"title,omitempty"`
	Body          *string        `json:"body,omitempty"`
	User          *github.User   `json:"user,omitempty"`
	Assignee      *github.User   `json:"assignee,omitempty"`
	ClosedAt      *time.Time     `json:"closed_at,omitempty"`
	CreatedAt     *time.Time     `json:"created_at,omitempty"`
	UpdatedAt     *time.Time     `json:"updated_at,omitempty"`
	ClosedBy      *github.User   `json:"closed_by,omitempty"`
	URL           *string        `json:"url,omitempty"`
	RepositoryURL *string        `json:"repository_url,omitempty"`
	Assignees     []*github.User `json:"assignees,omitempty"`
}

type IssueFilter struct {
	State string
}

type Commit struct {
	SHA       *string              `json:"sha,omitempty"`
	Author    *github.CommitAuthor `json:"author,omitempty"`
	Committer *github.CommitAuthor `json:"committer,omitempty"`
	Message   *string              `json:"message,omitempty"`
	URL       *string              `json:"url,omitempty"`
}

type CommitFilter struct {
	Author string
	Since  time.Time
	Until  time.Time
}
