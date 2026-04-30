package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/desulaidovich/plea-cli/generator"
	"github.com/desulaidovich/plea-cli/logger"
)

func main() {
	app := &cli.App{
		Name:        "plea-cli",
		Usage:       "Template generator for Go projects",
		Description: "Create new projects from ready-to-use templates",
		Version:     "dev",
		Authors: []*cli.Author{
			{
				Name:  "Anton Styazhkin",
				Email: "desulaidovich@icloud.com",
			},
		},
		Copyright: "(c) 2026",
		Commands: []*cli.Command{
			{
				Name:  "new",
				Usage: "Create a new project from template",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "Project name",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "repo",
						Aliases: []string{"r"},
						Usage:   "Repository path (e.g. github username or org)",
					},
					&cli.PathFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "Output directory",
					},
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "Show detailed output (commands, created files)",
					},
				},
				Action: func(ctx *cli.Context) error {
					loggen := logger.New(os.Stdout, logger.Debug, log.Ltime)
					gen := generator.New(ctx.String("name"), ctx.String("repo"), loggen, ctx.Bool("verbose"))
					return gen.Build()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "plea: %v\n", err)
		os.Exit(1)
	}
}
