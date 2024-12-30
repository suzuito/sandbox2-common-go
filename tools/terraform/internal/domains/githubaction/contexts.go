package githubaction

import (
	"strings"
)

type GitHubContext struct {
	EventName       string              `json:"event_name"`
	RefName         string              `json:"ref_name"`
	RefType         string              `json:"ref_type"`
	Repository      string              `json:"repository"`
	RepositoryOwner string              `json:"repository_owner"`
	Issue           *GithubContextIssue `json:"issue"`
	Event           GitHubContextEvent  `json:"event"`
}

func (t *GitHubContext) RepositoryName() string {
	return strings.Split(t.Repository, "/")[1]
}

type GitHubContextEvent struct {
	Comment *GitHubContextEventComment `json:"comment"`
}

type GitHubContextEventComment struct {
	Body string `json:"body"`
}

type GithubContextIssue struct {
	PullRequest *GithubContextIssuePullRequest `json:"pull_request"`
	Number      int                            `json:"number"`
}

type GithubContextIssuePullRequest struct {
}
