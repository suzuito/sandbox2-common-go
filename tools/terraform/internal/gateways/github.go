package gateways

import (
	"context"

	"github.com/google/go-github/v68/github"
)

type GithubPullRequestsService interface {
	Merge(ctx context.Context, owner string, repo string, number int, commitMessage string, opt *github.PullRequestOptions) (*github.PullRequestMergeResult, *github.Response, error)
	Get(ctx context.Context, owner string, repo string, number int) (*github.PullRequest, *github.Response, error)
	ListFiles(ctx context.Context, owner string, repo string, number int, opts *github.ListOptions) ([]*github.CommitFile, *github.Response, error)
}

type GithubIssuesService interface {
	CreateComment(ctx context.Context, owner string, repo string, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error)
}
