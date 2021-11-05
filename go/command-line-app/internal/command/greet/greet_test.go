package greet

import (
	"errors"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"command-line-app/internal/command"
	"command-line-app/internal/github"
)

func TestNew(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	assert.NotNil(t, c)
}

func TestNewFactory(t *testing.T) {
	ui := cli.NewMockUi()
	c, err := NewFactory(ui)()

	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestCommand_Synopsis(t *testing.T) {
	c := new(Command)
	synopsis := c.Synopsis()

	assert.NotEmpty(t, synopsis)
}

func TestCommand_Help(t *testing.T) {
	c := new(Command)
	help := c.Help()

	assert.NotEmpty(t, help)
}

func TestCommand_Run(t *testing.T) {
	t.Run("InvalidFlag", func(t *testing.T) {
		c := &Command{ui: cli.NewMockUi()}
		exitCode := c.Run([]string{"-undefined"})

		assert.Equal(t, command.FlagError, exitCode)
	})

	t.Run("OK", func(t *testing.T) {
		c := &Command{ui: cli.NewMockUi()}
		c.Run([]string{})

		assert.NotNil(t, c.services.github)
	})
}

func TestCommand_parseFlags(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedExitCode int
	}{
		{
			name:             "InvalidFlag",
			args:             []string{"-undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name:             "ValidFlag",
			args:             []string{"-username", "octocat"},
			expectedExitCode: command.Success,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Command{ui: cli.NewMockUi()}
			exitCode := c.parseFlags(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}

func TestCommand_exec(t *testing.T) {
	tests := []struct {
		name             string
		usernameFlag     string
		github           *MockGithubService
		expectedGreeting string
		expectedExitCode int
	}{
		{
			name:         "GetUserFails",
			usernameFlag: "octocat",
			github: &MockGithubService{
				GetUserMocks: []GetUserMock{
					{OutError: errors.New("http error")},
				},
			},
			expectedGreeting: "",
			expectedExitCode: command.GenericError,
		},
		{
			name:         "Success",
			usernameFlag: "octocat",
			github: &MockGithubService{
				GetUserMocks: []GetUserMock{
					{
						OutUser: &github.User{
							ID:    1,
							Login: "octocat",
							Email: "octocat@example.com",
							Name:  "Octocat",
						},
					},
				},
			},
			expectedGreeting: "Hello, Octocat!",
			expectedExitCode: command.Success,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Command{
				ui: cli.NewMockUi(),
			}

			c.flags.username = tc.usernameFlag
			c.services.github = tc.github

			exitCode := c.exec()

			assert.Equal(t, tc.expectedGreeting, c.Greeting())
			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
