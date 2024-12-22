package reporter

type Reporter interface {
	Reportf(path string, format string, args ...any)
}
