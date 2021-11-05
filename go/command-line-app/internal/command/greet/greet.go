package greet

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/mitchellh/cli"

	"command-line-app/internal/command"
	"command-line-app/internal/github"
)

const (
	timeout  = time.Minute
	synopsis = `Greet a GitHub user!`
	help     = `
  Use this command for greeting a GitHub user!

  Usage:  command-line-app greet [flags]

  Flags:
    -username  a GitHub username

  Examples:
    command-line-app greet -username octocat
    command-line-app greet -username=moorara
  `
)

type (
	githubService interface {
		GetUser(context.Context, string) (*github.User, error)
	}
)

// Command implements the cli.Command implementation.
type Command struct {
	ui    cli.Ui
	flags struct {
		username string
	}
	services struct {
		github githubService
	}
	outputs struct {
		greeting string
	}
}

// New creates a new command.
func New(ui cli.Ui) *Command {
	return &Command{
		ui: ui,
	}
}

// NewFactory returns a cli.CommandFactory for creating a new command.
func NewFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return New(ui), nil
	}
}

// Synopsis returns a short one-line synopsis for the command.
func (c *Command) Synopsis() string {
	return synopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *Command) Help() string {
	return help
}

// Run runs the actual command with the given command-line arguments.
// This method is used as a proxy for creating dependencies and the actual command execution is delegated to the run method for testing purposes.
func (c *Command) Run(args []string) int {
	if code := c.parseFlags(args); code != command.Success {
		return code
	}

	github, err := github.NewService()
	if err != nil {
		c.ui.Error(err.Error())
		return command.GenericError
	}

	c.services.github = github

	return c.exec()
}

func (c *Command) parseFlags(args []string) int {
	fs := flag.NewFlagSet("greet", flag.ContinueOnError)
	fs.StringVar(&c.flags.username, "username", "", "")

	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		// In case of error, the error and help will be printed by the Parse method
		return command.FlagError
	}

	return command.Success
}

// exec in an auxiliary method, so we can test the business logic with mock dependencies.
func (c *Command) exec() int {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if c.flags.username == "" {
		c.ui.Error("No GitHub username is provided.")
		return command.GenericError
	}

	user, err := c.services.github.GetUser(ctx, c.flags.username)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GenericError
	}

	c.outputs.greeting = fmt.Sprintf("Hello, %s!", user.Name)

	c.ui.Info(c.outputs.greeting)

	return command.Success
}

// Greeting returns the greeting text.
func (c *Command) Greeting() string {
	return c.outputs.greeting
}
