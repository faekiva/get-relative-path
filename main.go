package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func runApp(args ...string) (string, error) {
	var relativeTo string
	var err error
	output := ""
	app := &cli.App{
		Name:  "get-relative-path",
		Usage: "given a path, return it relative to the current working directory (or a specified path)",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "relative-to",
				Value: "",
				Usage: "path to use as the root of the relative path",
			},
		},
		Action: func(c *cli.Context) error {
			relativeTo = c.String("relative-to")
			if c.Value("relative-to") == "" {
				relativeTo, err = os.Getwd()
				if err != nil {
					return err
				}
			}
			if c.NArg() != 1 {
				return fmt.Errorf("expected 1 argument, got %d", c.Args().Len())
			}
			output, err = filepath.Rel(relativeTo, c.Args().First())
			if err != nil {
				return err
			}
			return nil
		},
	}
	err = app.Run(args)
	return output, err
}

func main() {
	output, err := runApp(os.Args...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
