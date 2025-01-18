package e2ehelpers

import (
	"errors"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PlaywrightTestCaseForSSR struct {
	Desc     string
	Setup    func(t *testing.T, testID TestID, exe *PlaywrightTestCaseForSSRExec)
	Teardown func(t *testing.T, testID TestID)
}

func (c *PlaywrightTestCaseForSSR) Run(t *testing.T) {
	if c.Setup == nil {
		panic(errors.New("setup function is required"))
	}

	testID := NewTestID()
	exe := PlaywrightTestCaseForSSRExec{}
	c.Setup(t, testID, &exe)

	pw, err := playwright.Run()
	require.NoError(t, err)
	defer pw.Stop() //nolint:errcheck

	browser, err := pw.Chromium.Launch()
	require.NoError(t, err)
	defer browser.Close()

	page, err := browser.NewPage()
	require.NoError(t, err)

	if exe.Do != nil {
		exe.Do(t, pw, browser, page)
	}
}

type PlaywrightTestCaseForSSRExec struct {
	Do func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page)
}

func AssertElementExists(t *testing.T, loc playwright.Locator) {
	c, err := loc.Count()
	require.NoError(t, err)
	assert.Greaterf(t, c, 0, "element %+v does not exist", loc)
}

func AssertElementNotExists(t *testing.T, loc playwright.Locator) {
	c, err := loc.Count()
	require.NoError(t, err)
	assert.Lessf(t, c, 1, "element %+v exists", loc)
}
