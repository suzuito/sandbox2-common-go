package terraformexe

import "fmt"

type PlanResult struct {
	IsPlanDiff bool
	Stdout     string
	Stderr     string
}

func (t *PlanResult) String() string {
	return fmt.Sprintf("out:\n%s\nerr:\n%s", t.Stdout, t.Stderr)
}

type ApplyResult struct {
	Stdout string
	Stderr string
}

func (t *ApplyResult) String() string {
	return fmt.Sprintf("out:\n%s\nerr:\n%s", t.Stdout, t.Stderr)
}
