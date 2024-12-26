package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var ErrPathIsntUseful = errors.New("path is not useful")

type CaseSensitivityGuesser func(paths ...string) bool
type MockableGetSysInfo func(path string) (any, error)

func getSysInfo(path string) (any, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return fmt.Sprint(info.Sys()), nil
}

func guessCaseSensitivity(paths ...string) bool {
	return guessCaseSensitivityInternal(runtime.GOOS, getSysInfo, paths...)
}

func guessCaseSensitivityInternal(operatingSystem string, stat MockableGetSysInfo, paths ...string) bool {
	if len(paths) == 0 {
		return getOSCaseSensitivityFallback(operatingSystem)
	}

	isThisPathCaseSensitive := func(path string) (isCaseSensitive bool, err error) {
		differentlyCasedPath, possible := getDifferentlyCasedVersionOfPath(path)
		if !possible {
			return false, ErrPathIsntUseful
		}
		info, err := stat(differentlyCasedPath)
		info2, err2 := stat(path)

		if err != nil && err2 != nil {
			// if they all error, then the path doesn't exist
			return false, err

		} else if err != nil || err2 != nil {
			// if 1-2 of them error, then it's case sensitive
			return true, nil

		} else if info == info2 {
			// if they all match, then it's case insensitive
			return false, nil
		}
		// if one of them points to different info, then it's case sensitive
		return true, nil
	}

	isAnyPathSensitive := false
	hasAnyPathReported := false
	for _, path := range paths {
		pathIsCaseSensitive, err := isThisPathCaseSensitive(path)
		if err != nil {
			continue
		}
		if pathIsCaseSensitive {
			isAnyPathSensitive = true
			hasAnyPathReported = true
			break
		}
		hasAnyPathReported = true
	}
	if hasAnyPathReported {
		return isAnyPathSensitive
	}

	return getOSCaseSensitivityFallback(operatingSystem)
}

func getDifferentlyCasedVersionOfPath(path string) (outPath string, possible bool) {
	lowerPath := strings.ToLower(path)
	pathIsSameAsLower := lowerPath == path

	if !pathIsSameAsLower {
		return lowerPath, true
	}

	upperPath := strings.ToUpper(path)
	pathIsSameAsUpper := upperPath == path

	if !pathIsSameAsUpper {
		return upperPath, true
	}

	return "", false
}

// If the answer is "under some common circumstances", then it's treated as case sensitive.
func getOSCaseSensitivityFallback(os string) bool {
	switch os {
	case "darwin":
		return false
	case "windows":
		return false
	default:
		return true
	}
}
