# command-line-app

This is intended to be used as a template for generating a new command-line application.

## API

| Command | Description |
|---------|-------------|
| `greet` | Creates and prints a greeting for a GitHub user! |

## Development

### Make

| Rule | Description |
|------|-------------|
| `test` | Runs the unit tests with `-race` flag. |
| `test-short` | Runs the unit tests with `-short` flag. |
| `test-coverage` | Runs the unit tests and generates coverage reports (`c.out` and `coverage.html`). |
| `clean-test` | Deletes files generated by tests. |
| `run` | Runs the application. |
| `build` | Builds the application binary. |
| `build-all` | Builds the application binary for all supported platforms. |
| `clean-build` | Deletes built binaries. |
