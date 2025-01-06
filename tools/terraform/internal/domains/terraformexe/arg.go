package terraformexe

import (
	"encoding/json"
	"os"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/githubaction"
)

type Arg struct {
	TargetType              TargetType
	PlanOnly                bool
	GitHubOwner             string
	GitHubRepository        string
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
	eventName string,
	eventPath string,
) (*Arg, bool, error) {
	eventPayloadBytes, err := os.ReadFile(eventPath)
	if err != nil {
		return nil, false, terrors.Wrap(err)
	}

	switch eventName {
	case "issue_comment":
		eventPayload := githubaction.PayloadIssueComment{}
		if err := json.Unmarshal(eventPayloadBytes, &eventPayload); err != nil {
			return nil, false, terrors.Wrap(err)
		}

		arg := Arg{
			TargetType:              ForOnlyChageFiles,
			GitHubOwner:             eventPayload.Repository.Owner.Name,
			GitHubRepository:        eventPayload.Repository.Name,
			GitHubPullRequestNumber: eventPayload.Issue.Number,
		}
		switch eventPayload.Comment.Body {
		case "///terraform plan":
			arg.PlanOnly = true
		case "///terraform apply":
			arg.PlanOnly = false
		default:
			return nil, false, nil
		}

		return &arg, true, nil
	case "schedule", "workflow_dispatch":
		eventPayload := githubaction.PayloadSchedule{}
		if err := json.Unmarshal(eventPayloadBytes, &eventPayload); err != nil {
			return nil, false, terrors.Wrap(err)
		}

		arg := Arg{
			TargetType:       ForAllFiles,
			PlanOnly:         true,
			GitHubOwner:      eventPayload.Repository.Owner.Name,
			GitHubRepository: eventPayload.Repository.Name,
		}

		return &arg, true, nil
	}

	return nil, false, nil
}
