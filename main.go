package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
)

func runApp(args Args, guessCaseSensitive CaseSensitivityGuesser) (string, error) {
	path := args.Path
	relativeTo := args.RelativeTo
	var err error
	var isCaseSensitive bool

	switch args.IsCaseSensitive {
	case "true":
		isCaseSensitive = true
	case "false":
		isCaseSensitive = false
	default:
		isCaseSensitive = guessCaseSensitive(relativeTo, path)
	}

	if !filepath.IsAbs(relativeTo) {
		relativeTo, err = filepath.Abs(relativeTo)
		if err != nil {
			return "", err
		}
	}

	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}

	if !isCaseSensitive {
		relativeTo = strings.ToLower(relativeTo)
		path = strings.ToLower(path)
	}

	output, err := filepath.Rel(relativeTo, path)
	return output, err
}

type Args struct {
	RelativeTo      string `arg:"--relative-to" default:"."`
	Path            string `arg:"positional" help:"if provided path is relative, it will be resolved relative to PWD first, then relative to the path provided with --relative-to"`
	IsCaseSensitive string `arg:"-c, --case-sensitive" default:"guess" help:"options are true, false, or guess"`
}

type RunAppArgs struct {
	Args
	IsCaseSensitive bool
}

func main() {
	var args Args
	arg.MustParse(&args)
	output, err := runApp(args, guessCaseSensitivity)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
