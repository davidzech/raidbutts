package simc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	const version = Version(`SimulationCraft 902-01 for World of Warcraft 9.0.2.37176 Live (hotfix 2021-01-17/37176, git build shadowlands d353c77f0e)`)

	assert.Equal(t, "902-01", version.Release())
	assert.Equal(t, "shadowlands", version.Branch())
	assert.Equal(t, "d353c77f0e", version.GitHash())
	assert.Equal(t, "9.0.2.37176", version.WoWVersion())
	assert.Equal(t, "Live", version.WoWReleaseType())
	assert.Equal(t, "2021-01-17", version.WoWHotfixDate())
	assert.Equal(t, "37176", version.WoWHotfix())
}
