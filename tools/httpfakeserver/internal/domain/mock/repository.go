package mock

import (
	"iter"
	"maps"
	"slices"
	"sync"
)

type Repository struct {
	casesMu sync.Mutex
	cases   map[string]Mock
}

func (m *Repository) SetMock(c Mock) {
	m.casesMu.Lock()
	defer m.casesMu.Unlock()
	m.cases[c.ID()] = c
}

func (m *Repository) Clear() {
	m.casesMu.Lock()
	defer m.casesMu.Unlock()
	m.cases = map[string]Mock{}
}

func (m *Repository) Mocks() iter.Seq[Mock] {
	return func(yield func(Mock) bool) {
		m.casesMu.Lock()
		defer m.casesMu.Unlock()
		for _, key := range slices.Sorted(maps.Keys(m.cases)) {
			if !yield(m.cases[key]) {
				break
			}
		}
	}
}

func NewRepository() *Repository {
	return &Repository{
		cases: map[string]Mock{},
	}
}
