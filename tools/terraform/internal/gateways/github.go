package gateways

import (
	"context"

	"github.com/google/go-github/v68/github"
)

type GithubPullRequestsService interface {
	ListFiles(ctx context.Context, owner string, repo string, number int, opts *github.ListOptions) ([]*github.CommitFile, *github.Response, error)
}
