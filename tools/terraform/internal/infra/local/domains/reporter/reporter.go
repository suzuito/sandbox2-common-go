package reporter

import "fmt"

type impl struct{}

func (t *impl) Reportf(path string, format string, args ...any) {
	fmt.Printf(
		"%s %s\n",
		fmt.Sprintf(format, args...),
		path,
	)
}

func New() *impl {
	return &impl{}
}
