package terraformexe

import (
	"encoding/json"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/githubaction"
)

type Arg struct {
	TargetType              TargetType
	PlanOnly                bool
	GitHubOwner             string
	GitHubRepository        string
	GitHubRefName           string
	GitHubRefType           string
	GitHubPullRequestNumber int
}

type TargetType int

const (
	ForOnlyChageFiles TargetType = iota + 1
	ForAllFiles
)

func NewTerraformExecutionArg(
	dirPathBase string,
	projectID string,
	githubContextJSON string,
) (*Arg, bool, error) {
	githubContext := githubaction.GitHubContext{}
	if err := json.Unmarshal([]byte(githubContextJSON), &githubContext); err != nil {
		return nil, false, terrors.Errorf("github context is invalid json: %s: %w", githubContextJSON, err)
	}

	if githubContext.EventName == "issue_comment" &&
		githubContext.Issue.PullRequest != nil &&
		githubContext.Event.Comment != nil &&
		githubContext.Event.Comment.Body == "///terraform plan" {
		// comment on pull request
		if githubContext.Issue == nil {
			return nil, false, terrors.Errorf("invalid github context: issue is null")
		}

		return &Arg{
			TargetType:              ForOnlyChageFiles,
			PlanOnly:                true,
			GitHubOwner:             githubContext.RepositoryOwner,
			GitHubRepository:        githubContext.RepositoryName(),
			GitHubPullRequestNumber: githubContext.Issue.Number,
		}, true, nil
	}

	if githubContext.EventName == "issue_comment" &&
		githubContext.Issue.PullRequest != nil &&
		githubContext.Event.Comment != nil &&
		githubContext.Event.Comment.Body == "///terraform apply" {
		// comment on pull request
		if githubContext.Issue == nil {
			return nil, false, terrors.Errorf("invalid github context: issue is null")
		}

		return &Arg{
			TargetType:              ForOnlyChageFiles,
			GitHubOwner:             githubContext.RepositoryOwner,
			GitHubRepository:        githubContext.RepositoryName(),
			GitHubPullRequestNumber: githubContext.Issue.Number,
		}, true, nil
	}

	if githubContext.EventName == "schedule" {
		return &Arg{
			TargetType:       ForAllFiles,
			PlanOnly:         true,
			GitHubOwner:      githubContext.RepositoryOwner,
			GitHubRepository: githubContext.RepositoryName(),
			GitHubRefName:    githubContext.RefName,
			GitHubRefType:    githubContext.RefType,
		}, true, nil
	}

	return nil, false, nil
}
