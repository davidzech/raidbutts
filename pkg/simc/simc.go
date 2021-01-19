package simc

import (
	"regexp"
)

type Configuration string

type Result struct {
	Data map[string]interface{}
}

type Simulator interface {
	Simulate(Configuration) (*Result, error)
	Version() (Version, error)
}

type Version string

const versionRegex = `SimulationCraft (.+) for World of Warcraft (.+) (.+) \(hotfix (.+)/(.+), git build (.+) (.+)\)`
const versionRegexGroups = 7

var versionRegexp = regexp.MustCompile(versionRegex)

func (v Version) Release() string {
	return v.regexCaptureGroup(1)
}

func (v Version) WoWVersion() string {
	return v.regexCaptureGroup(2)
}

func (v Version) WoWReleaseType() string {
	return v.regexCaptureGroup(3)
}

func (v Version) WoWHotfixDate() string {
	return v.regexCaptureGroup(4)
}

func (v Version) WoWHotfix() string {
	return v.regexCaptureGroup(5)
}

func (v Version) Branch() string {
	return v.regexCaptureGroup(6)
}

func (v Version) GitHash() string {
	return v.regexCaptureGroup(7)
}

func (v Version) regexCaptureGroup(n int) string {
	matches := versionRegexp.FindStringSubmatch(string(v))
	if len(matches) != versionRegexGroups+1 {
		return ""
	}
	return matches[n]
}
