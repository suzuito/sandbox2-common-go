package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueNamePublicServices(t *testing.T) {
	countsServiceID := map[ServiceID]int{}
	countsUniqueName := map[string]int{}

	for _, s := range availableServices {
		countServiceID := countsServiceID[s.ID]
		countServiceID += 1
		countsServiceID[s.ID] = countServiceID

		countUniqueName := countsUniqueName[s.UniqueName]
		countUniqueName += 1
		countsUniqueName[s.UniqueName] = countUniqueName
	}

	for serviceID, count := range countsServiceID {
		assert.LessOrEqual(
			t,
			count,
			1,
			fmt.Sprintf(
				"%s is duplicated",
				serviceID.UUID(),
			),
		)
	}

	for uniqueName, count := range countsUniqueName {
		assert.LessOrEqual(
			t,
			count,
			1,
			fmt.Sprintf(
				"%s is duplicated",
				uniqueName,
			),
		)
	}
}
