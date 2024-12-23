package reporter

import "fmt"

type impl struct{}

func (t *impl) Reportf(path string, format string, args ...any) {
	fmt.Printf(
		"%s (%s)\n",
		fmt.Sprintf(format, args...),
		path,
	)
}

func (t *impl) AssertEqualf(path string, expected, actual any, format string, args ...any) bool {
	eq := expected == actual

	if !eq {
		fmt.Printf(
			"%s (%s)\n",
			fmt.Sprintf(format, args...),
			path,
		)
		fmt.Printf("  expected: %s\n", expected)
		fmt.Printf("  actual: %s\n", actual)
	}

	return eq
}

func (t *impl) AssertTruef(path string, actual bool, format string, args ...any) bool {
	if !actual {
		fmt.Printf(
			"%s (%s)\n",
			fmt.Sprintf(format, args...),
			path,
		)
		fmt.Println("  expected: true")
		fmt.Printf("  actual: %v\n", actual)
	}

	return true
}

func New() *impl {
	return &impl{}
}
