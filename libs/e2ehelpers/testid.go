package e2ehelpers

import "github.com/google/uuid"

type TestID uuid.UUID

func (t *TestID) UUID() uuid.UUID {
	return uuid.UUID(*t)
}

func (t *TestID) String() string {
	return t.UUID().String()
}

func NewTestID() TestID {
	return TestID(uuid.New())
}
