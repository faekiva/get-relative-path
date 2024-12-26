package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const homeDir = "/Users/faekiva"

func runCaseSensitive(t *testing.T, args ...string) (string, error) {
	t.Helper()

	os.Args = append([]string{"get-relative-path", "--case-sensitive", "true"}, args...)
	return runApp(guessCaseSensitive)
}

func runCaseInsensitive(t *testing.T, args ...string) (string, error) {
	t.Helper()
	os.Args = append([]string{"get-relative-path", "--case-sensitive", "false"}, args...)
	return runApp(guessCaseSensitive)
}

func runWithGuesser(t *testing.T, guesser CaseSensitivityGuesser, args ...string) (string, error) {
	t.Helper()
	os.Args = append([]string{"get-relative-path"}, args...)
	return runApp(guesser)
}

func TestTwoArgsNoFlag(t *testing.T) {
	_, err := runCaseSensitive(t, "/Users/faekiva/go/src/github.com/kiva/get-relative-path", "/Users/faekiva/go/src/github.com/kiva/get-relative-path/main.go")
	assert.Error(t, err)
}

func TestSamePath(t *testing.T) {
	output, err := runCaseSensitive(t, homeDir, "--relative-to", homeDir+"/")
	require.NoError(t, err)
	assert.Equal(t, ".", output)
}

func TestChildPath(t *testing.T) {
	output, err := runCaseSensitive(t, homeDir, "--relative-to", "/Users")
	require.NoError(t, err)
	assert.Equal(t, "faekiva", output)
}

func TestCaseInsensitive(t *testing.T) {
	output, err := runCaseInsensitive(t, homeDir, "--relative-to", "/users")
	require.NoError(t, err)
	assert.Equal(t, "faekiva", output)
}

func TestCaseInsensitiveGuesser(t *testing.T) {
	guesser := func(paths ...string) bool {
		return false
	}
	output, err := runWithGuesser(t, guesser, homeDir, "--relative-to", "/users")
	require.NoError(t, err)
	assert.Equal(t, "faekiva", output)
}
