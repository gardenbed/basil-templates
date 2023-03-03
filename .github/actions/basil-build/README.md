# Basil Build

This action is used for building Go binaries using Basil CLI.

## Inputs

### `path`

The path to the Go project. The default value is `.`.

By convention, Basil assumes the project directory is a main package if it contains a main.go file.
It also assumes every directory inside `cmd` is a main package for a binary with the same name as the directory name.

### `cross_compile`

If set to `true`, the binaries will be built for all platforms.
This will override any value set in the project `basil.yaml` file.

### `platforms`

A comma-separated list of platforms for cross compilation.
This will override any value set in the project `basil.yaml` file.

## Outputs

### `artifacts`

A space-separated list of built artifacts (binary files).

## Example Usages

```yaml
name: Example
on: [push]
jobs:
  build:
    name: Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - name: Basil Build
        uses: ./.github/actions/basil-build
```

```yaml
name: Example
on: [push]
jobs:
  build:
    name: Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - name: Basil Build
        uses: ./.github/actions/basil-build
        with:
          path: ./services/hello-world
```

```yaml
name: Example
on: [push]
jobs:
  build:
    name: Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - name: Basil Build
        uses: ./.github/actions/basil-build
        with:
          cross_compile: 'true'
          platforms: 'linux-amd64,linux-arm64,darwin-amd64,darwin-arm64,windows-amd64,windows-arm64'
```
