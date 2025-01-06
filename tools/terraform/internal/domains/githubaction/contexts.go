package githubaction

type PayloadIssueComment struct {
	Comment    *Comment    `json:"comment"`
	Issue      *Issue      `json:"issue"`
	Repository *Repository `json:"repository"`
}

type Comment struct {
	Body string `json:"body"`
}

type Issue struct {
	PullRequest *PullRequest `json:"pull_request"`
	Number      int          `json:"number"`
}

type PullRequest struct {
}

type PayloadSchedule struct {
	Ref        string      `json:"ref"`
	Repository *Repository `json:"repository"`
}

type Repository struct {
	Name  string          `json:"name"`
	Owner RepositoryOwner `json:"owner"`
}

type RepositoryOwner struct {
	Name string `json:"name"`
}
