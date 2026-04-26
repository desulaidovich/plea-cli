package main

import (
	"fmt"
	"log"
	"os"
	"plea-cli/internal/generator"
	"plea-cli/internal/logger"

	"github.com/urfave/cli/v2"
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
						Usage:    "Mmodule name in go.mod",
						Value:    ".",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "Output directory",
						Value:    ".",
						Required: true,
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
					projectName := ctx.String("name")
					moduleName := ctx.String("module")
					outputDir := ctx.String("output")
					logLevel := ctx.String("log-level")
					verbose := ctx.Bool("verbose")
					loggen := logger.New(os.Stdout, logger.Debug, log.Ltime)
					return generator.New(projectName, moduleName, outputDir, logLevel, verbose, loggen).Do()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
