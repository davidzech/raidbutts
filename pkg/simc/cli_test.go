package simc

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkSimc(t *testing.T) {
	_, err := exec.LookPath("simc")
	if err != nil {
		t.Skip("simc not installed")
	}
}

func TestCLI_Simulate(t *testing.T) {
	checkSimc(t)
	cli := CLI{Executable: "simc"}
	confBytes, err := ioutil.ReadFile("testdata/profile.simc")
	require.NoError(t, err)
	testProfile := Configuration(string(confBytes))

	result, err := cli.Simulate(testProfile)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestCLI_Version(t *testing.T) {
	checkSimc(t)

	cli := CLI{Executable: "simc"}

	ver, err := cli.Version()
	require.NoError(t, err)
	r := regexp.MustCompile(versionRegex)
	assert.Regexp(t, r, ver)
	assert.NotEmpty(t, ver.GitHash())
}
