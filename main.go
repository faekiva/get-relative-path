package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	arg "github.com/alexflint/go-arg"
)

func runApp(guessCaseSensitive CaseSensitivityGuesser) (string, error) {
	args := Args{}
	err := arg.Parse(&args)
	if err != nil {
		return "", err
	}
	path := args.Path
	relativeTo := args.RelativeTo

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
	IsCaseSensitive string `arg:"-c, --case-sensitive" default:"guess"`
}

type RunAppArgs struct {
	Args
	IsCaseSensitive bool
}

func main() {
	output, err := runApp(guessCaseSensitivity)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
