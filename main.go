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

	var isCaseSensitive bool

	switch args.IsCaseSensitive {
	case "true":
		isCaseSensitive = true
	case "false":
		isCaseSensitive = false
	default:
		isCaseSensitive = guessCaseSensitive(args.RelativeTo, args.Path)
	}

	if !isCaseSensitive {
		args.RelativeTo = strings.ToLower(args.RelativeTo)
		args.Path = strings.ToLower(args.Path)
	}

	output, err := filepath.Rel(args.RelativeTo, args.Path)
	return output, err
}

type Args struct {
	RelativeTo      string `arg:"--relative-to" default:"."`
	Path            string `arg:"positional"`
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
