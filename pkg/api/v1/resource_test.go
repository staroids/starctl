package v1

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type parseProjectTestStruct struct {
	flag   string
	parsed []string
}

var parseProjectTestData = []parseProjectTestStruct{
	{"GITHUB/staroid/app:master#commit1", []string{"GITHUB", "staroid", "app", "master", "commit1", "false"}},
	{"GITHUB/staroid/app:master", []string{"GITHUB", "staroid", "app", "master", "", "false"}},
	{"GITHUB/staroid/app", []string{"", "", "", "", "", "true"}},
	{"GITHUB/staroid/app:", []string{"", "", "", "", "", "true"}},
	{"GITHUB/staroid:master", []string{"", "", "", "", "", "true"}},
}

func TestNewCommitFromCommitLocation(t *testing.T) {
	for _, testData := range parseProjectTestData {
		commit, err := NewCommitFromCommitLocation(testData.flag)
		if err == nil {
			assert.Equal(t, testData.parsed[0], commit.Provider)
			assert.Equal(t, testData.parsed[1], commit.Owner)
			assert.Equal(t, testData.parsed[2], commit.Repo)
			assert.Equal(t, testData.parsed[3], commit.Branch)
			assert.Equal(t, testData.parsed[4], commit.Commit)
		}
		assert.Equal(t, testData.parsed[5], strconv.FormatBool(err != nil))
	}
}
