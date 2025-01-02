package domains

type State struct {
	ExecutedHistories ExecutedHistories
}

type ExecutedHistories []ExecutedHistory

type ExecutedHistory struct {
	Args []string
}
