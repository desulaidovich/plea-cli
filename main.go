package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/desulaidovich/plea-cli/internal/generator"
	"github.com/desulaidovich/plea-cli/internal/logger"
	"github.com/desulaidovich/plea-cli/internal/manifest"
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
						Name:     "module",
						Aliases:  []string{"m"},
						Usage:    "Module name in go.mod",
						Value:    ".",
						Required: true,
					},
					&cli.PathFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "Output directory",
						Value:   ".",
					},
					&cli.StringFlag{
						Name:    "log-level",
						Aliases: []string{"ll"},
						Usage:   "Log level for slog (debug, info, warn, error)",
						Value:   "debug",
						Action: func(ctx *cli.Context, level string) error {
							allowed := map[string]bool{
								"debug": true,
								"info":  true,
								"warn":  true,
								"error": true,
							}

							if !allowed[level] {
								return fmt.Errorf("log-level must be one of: debug, info, warn, error")
							}
							return nil
						},
					},
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "Show detailed output (commands, created files)",
					},
				},
				Action: func(ctx *cli.Context) error {
					cfg := manifest.Config{
						Name:     ctx.String("name"),
						Module:   ctx.String("module"),
						Output:   ctx.String("output"),
						LogLevel: ctx.String("log-level"),
						Verbose:  ctx.Bool("verbose"),
					}
					loggen := logger.New(os.Stdout, logger.Debug, log.Ltime)
					return generator.New(cfg, loggen).Do()
				},
			},
			{
				Name:  "manifest",
				Usage: "Create a new project from template by plea.yaml",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:    "path",
						Aliases: []string{"p"},
						Usage:   "Path to plea.yaml",
						Value:   "./plea.yaml",
					},
				},
				Action: func(ctx *cli.Context) error {
					path := ctx.String("path")

					cfg, err := manifest.ReadFile(path)
					if err != nil {
						return fmt.Errorf("read manifest file: %w", err)
					}
					loggen := logger.New(os.Stdout, logger.Debug, log.Ltime)
					return generator.New(*cfg, loggen).Do()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "plea: %v\n", err)
		os.Exit(1)
	}
}
