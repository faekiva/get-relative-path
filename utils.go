package main

import (
	"time"
)

func getTmpFileName() string {
	return "get-relative-path" + time.Now().String()
}

func guessCaseSensitive(relativeTo, path string) bool {
	// guess := func(path string) bool {
	// 	os.Stat(path)
	// }
	return true
}
