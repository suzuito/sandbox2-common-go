package domains

import "fmt"

type IncrementType string

func (t *IncrementType) String() string {
	return string(*t)
}

func (t *IncrementType) Validate() error {
	switch *t {
	case IncrementTypeMajor, IncrementTypeMinor, IncrementTypePatch:
		return nil
	}
	return fmt.Errorf("unknown increment type %s", *t)
}

const (
	IncrementTypeMajor IncrementType = "major"
	IncrementTypeMinor IncrementType = "minor"
	IncrementTypePatch IncrementType = "patch"
)
