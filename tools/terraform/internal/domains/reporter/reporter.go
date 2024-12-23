package reporter

type Reporter interface {
	Reportf(path string, format string, args ...any)
	AssertEqualf(path string, expected, actual any, format string, args ...any) bool
	AssertTruef(path string, actual bool, format string, args ...any) bool
}
