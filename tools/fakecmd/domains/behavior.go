package domains

type BehaviorType int

const (
	BehaviorTypeStdoutStderrExitCode = iota + 1
)

type Behaviors []Behavior

type Behavior struct {
	Type                         BehaviorType
	BehaviorStdoutStderrExitCode *BehaviorStdoutStderrExitCode
}

type BehaviorStdoutStderrExitCode struct {
	Stdout   string
	Stderr   string
	ExitCode uint8
}
