package main

import (
	"os"
	"testing"

	"github.com/alexflint/go-arg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const homeDir = "/Users/faekiva"

func runCaseSensitive(t *testing.T, args ...string) (string, error) {
	t.Helper()
	os.Args = append([]string{"get-relative-path", "--case-sensitive", "true"}, args...)
	var cliArgs Args
	err := arg.Parse(&cliArgs)
	require.NoError(t, err)
	return runApp(cliArgs, guessCaseSensitivity)
}

func runCaseInsensitive(t *testing.T, args ...string) (string, error) {
	t.Helper()
	os.Args = append([]string{"get-relative-path", "--case-sensitive", "false"}, args...)
	var cliArgs Args
	err := arg.Parse(&cliArgs)
	require.NoError(t, err)
	return runApp(cliArgs, guessCaseSensitivity)
}

func runWithGuesser(t *testing.T, guesser CaseSensitivityGuesser, args ...string) (string, error) {
	t.Helper()
	os.Args = append([]string{"get-relative-path"}, args...)
	var cliArgs Args
	err := arg.Parse(&cliArgs)
	require.NoError(t, err)
	return runApp(cliArgs, guesser)
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

func TestAlwaysStartWithDot(t *testing.T) {
	output, err := runCaseInsensitive(t, homeDir, "--always-start-with-dot", "--relative-to", "/users")
	require.NoError(t, err)
	assert.Equal(t, "./faekiva", output)
}

func TestAbsolutePathAgainstPeriod(t *testing.T) {
	tmpDir := os.TempDir()
	testDir, err := os.MkdirTemp(tmpDir, "boop")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	testDir2, err := os.MkdirTemp(testDir, "beep")
	require.NoError(t, err)
	defer os.RemoveAll(testDir2)

	initialWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(initialWd)
	err = os.Chdir(testDir)

	require.NoError(t, err)

	// There can technically be multiple relative paths to the same directory because of symlinks,
	// Setting the PWD env var helps us get the most predicable result
	t.Setenv("PWD", testDir)
	where, err := runCaseInsensitive(t, tmpDir, "--relative-to", ".")
	require.NoError(t, err)
	assert.Equal(t, "..", where)

	where, err = runCaseInsensitive(t, "..", "--relative-to", testDir2)
	require.NoError(t, err)
	assert.Equal(t, "../..", where)

	t.Log("tmpDir:", tmpDir)
	t.Log("testDir:", testDir)
	t.Log("testDir2:", testDir2)
}
