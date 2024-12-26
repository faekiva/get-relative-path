package main

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testPath = "/Users/faekiva"
const testPathDifferentlyCased = "/Users/Faekiva"
const testPathLinuxNFS = "/volume1/linuxNFS"
const testPathWindowsSMB = "/volume1/windowsSMB"
const testPathNoAlphabet = "/2/12345"
const testLowerCasePath = "/users/faekiva"

var ErrAlreadyCalled = errors.New("already called")
var ErrPathDoesntExist = errors.New("path doesn't exist")

func osStatCaseInsensitiveEveryPathExists(path string) (any, error) {
	return 3, nil
}

func osStatCaseSensitiveEveryPathExists(path string) (any, error) {
	return path, nil
}

func noPathResolves(path string) (any, error) {
	return nil, ErrPathDoesntExist
}

func osStatPathExistsInOneCase() MockableGetSysInfo {
	previousCallsMap := make(map[string]string)
	return func(path string) (any, error) {
		lower := strings.ToLower(path)
		lastCalledWith, ok := previousCallsMap[lower]
		if ok {
			if path == lastCalledWith {
				return path, nil
			}
			return nil, ErrAlreadyCalled
		}
		previousCallsMap[lower] = path
		return path, nil
	}
}

func osStatOnePathSensitiveOnePathNot(sensitivePath, nonSensitivePath string) MockableGetSysInfo {
	nonSensitivePathLower := strings.ToLower(nonSensitivePath)
	return func(path string) (any, error) {
		if path == sensitivePath {
			return path, nil
		}
		lower := strings.ToLower(path)
		if lower == nonSensitivePathLower {
			return nonSensitivePath, nil
		}
		return path, nil
	}
}

func TestOsStatPathExistsInOneCase(t *testing.T) {
	testFunc := osStatPathExistsInOneCase()
	_, err := testFunc(testPath)
	require.NoError(t, err)
	_, err = testFunc(testPath)
	require.NoError(t, err)
	_, err = testFunc(testPathDifferentlyCased)
	require.Error(t, err)
}

// GIVEN: the OS is case insensitive
// AND the path exists in one case
// AND the path isn't on an case sensitive path
// WHEN I call guessCaseSensitive
// THEN: it should return false
func TestAllInsensitive(t *testing.T) {
	isCaseSensitive := guessCaseSensitivityInternal("darwin", osStatCaseInsensitiveEveryPathExists, testPath)
	require.False(t, isCaseSensitive)
}

// GIVEN: the OS is case insensitive
// AND the path exists in one case
// AND the path is on an case sensitive path
// WHEN I call guessCaseSensitive
// THEN: it should return false
func TestOnlyPathSensitive(t *testing.T) {
	isCaseSensitive := guessCaseSensitivityInternal("darwin", osStatCaseSensitiveEveryPathExists, testPath)
	require.True(t, isCaseSensitive)
}

// GIVEN: the path doesn't resolve
// AND the OS is case insensitive
// WHEN I call guessCaseSensitive
// THEN: it should return false
func TestNoPathResolvesInsensitive(t *testing.T) {
	isCaseSensitive := guessCaseSensitivityInternal("darwin", noPathResolves, testPath)
	require.False(t, isCaseSensitive)
}

// GIVEN: the path doesn't resolve
// AND the OS is case sensitive
// WHEN I call guessCaseSensitive
// THEN: it should return false
func TestNoPathResolvesSensitiveButOSIs(t *testing.T) {
	isCaseSensitive := guessCaseSensitivityInternal("linux", noPathResolves, testLowerCasePath)
	require.True(t, isCaseSensitive)
}

// GIVEN: the OS is case sensitive
// AND the path exists in only one case
// WHEN I call guessCaseSensitive
// THEN: it should return true
func TestAllSensitive(t *testing.T) {
	isCaseSensitive := guessCaseSensitivityInternal("linux", osStatPathExistsInOneCase(), testPath)
	require.True(t, isCaseSensitive)
}

// GIVEN: the OS is case sensitive
// AND the path is the same system object in all cases
// WHEN I call guessCaseSensitive
// THEN: it should return false
func TestAllSameCase(t *testing.T) {
	isCaseSensitive := guessCaseSensitivityInternal("linux", osStatCaseInsensitiveEveryPathExists, testPath)
	require.False(t, isCaseSensitive)
}

// GIVEN: the OS is anything
// AND one path is the same system object in all cases
// AND the other path is a different system object in all cases
// WHEN I call guessCaseSensitive
// THEN: it should return true
func TestOneSameCaseOneDifferentCase(t *testing.T) {
	for _, os := range []string{"darwin", "linux"} {
		isCaseSensitive := guessCaseSensitivityInternal(os, osStatOnePathSensitiveOnePathNot(testPathLinuxNFS, testPathWindowsSMB), testPathLinuxNFS, testPathWindowsSMB)
		require.True(t, isCaseSensitive)
		isCaseSensitive = guessCaseSensitivityInternal(os, osStatOnePathSensitiveOnePathNot(testPathLinuxNFS, testPathWindowsSMB), testPathWindowsSMB, testPathLinuxNFS)
		require.True(t, isCaseSensitive)
	}
}

func TestGivenNoPathsAreUseful(t *testing.T) {
	for _, os := range []string{"darwin", "linux"} {
		isCaseSensitive := guessCaseSensitivityInternal(os, osStatCaseInsensitiveEveryPathExists, testPathNoAlphabet)
		require.Equal(t, getOSCaseSensitivityFallback(os), isCaseSensitive)
	}
}
func TestGivenNoPathsAreProvided(t *testing.T) {
	for _, os := range []string{"darwin", "linux"} {
		isCaseSensitive := guessCaseSensitivityInternal(os, osStatCaseInsensitiveEveryPathExists)
		require.Equal(t, getOSCaseSensitivityFallback(os), isCaseSensitive)
	}
}

func TestWithCaseInsensitiveSys(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "asd")
	defer os.RemoveAll(tmpDir)
	statter := func(path string) (any, error) {
		info, err := getSysInfo(tmpDir)
		require.NoError(t, err)
		return info, nil
	}
	require.NoError(t, err)

	isCaseSensitive := guessCaseSensitivityInternal("darwin", statter, tmpDir, strings.ToUpper(tmpDir))
	require.False(t, isCaseSensitive)
}
