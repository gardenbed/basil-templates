package main

import (
	"os"

	"github.com/mitchellh/cli"

	"command-line-app/internal/command/greet"
	"command-line-app/metadata"
)

func main() {
	ui := createUI()
	app := createCLI(ui)

	code, err := app.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(code)
}

func createUI() cli.Ui {
	return &cli.ConcurrentUi{
		Ui: &cli.ColoredUi{
			Ui: &cli.BasicUi{
				Reader:      os.Stdin,
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
			OutputColor: cli.UiColorNone,
			InfoColor:   cli.UiColorGreen,
			WarnColor:   cli.UiColorYellow,
			ErrorColor:  cli.UiColorRed,
		},
	}
}

func createCLI(ui cli.Ui) *cli.CLI {
	c := cli.NewCLI("command-line-app", metadata.String())
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"greet": greet.NewFactory(ui),
	}

	return c
}
