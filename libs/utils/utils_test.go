package utils_test

import (
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
)

func TestPtr(t *testing.T) {
	utils.Ptr(1)
}

func TestFilter(t *testing.T) {
	testCases := []struct {
		desc     string
		input    []int
		inputF   func(int) bool
		expected []int
	}{
		{
			desc:   "empty",
			inputF: func(_ int) bool { return true },
		},
		{
			desc:     "not empty",
			input:    []int{1, 2, 3, 1, 2, 3},
			inputF:   func(i int) bool { return i == 2 },
			expected: []int{2, 2},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := slices.Collect(utils.Filter(tC.inputF, slices.Values(tC.input)))
			require.Equal(t, tC.expected, actual)
		})
	}

	// breakしたときにpanicしないこと(=yield関数がfalseを返した後で呼ばれていないこと)
	for range utils.Filter(
		func(_ int) bool { return true },
		slices.Values([]int{1, 2, 3}),
	) {
		break
	}
}

func TestMap(t *testing.T) {
	testCases := []struct {
		desc     string
		input    []int
		inputF   func(int) string
		expected []string
	}{
		{
			desc:   "empty",
			inputF: func(_ int) string { return "" },
		},
		{
			desc:     "not empty",
			input:    []int{1, 2, 3},
			inputF:   func(i int) string { return strconv.Itoa(i) },
			expected: []string{"1", "2", "3"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := slices.Collect(utils.Map(tC.inputF, slices.Values(tC.input)))
			require.Equal(t, tC.expected, actual)
		})
	}

	// breakしたときにpanicしないこと(=yield関数がfalseを返した後で呼ばれていないこと)
	for range utils.Map(
		func(_ int) string { return "" },
		slices.Values([]int{1, 2, 3}),
	) {
		break
	}
}
