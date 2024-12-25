package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	args = append([]string{"get-relative-path"}, args...)
	return runApp(args...)
}
func TestTwoArgsNoFlag(t *testing.T) {
	_, err := run(t, "/Users/faekiva/go/src/github.com/kiva/get-relative-path", "/Users/faekiva/go/src/github.com/kiva/get-relative-path/main.go")
	assert.Error(t, err)
}

func TestSamePath(t *testing.T) {
	output, err := run(t, "/Users/kiva", "--relative-to", "/Users/kiva/")
	require.NoError(t, err)
	assert.Equal(t, ".", output)
}

func TestChildPath(t *testing.T) {
	output, err := run(t, "/Users/kiva", "--relative-to", "/Users")
	require.NoError(t, err)
	assert.Equal(t, ".", output)
}
