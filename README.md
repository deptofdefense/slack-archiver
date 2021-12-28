# slack-archiver

## Description

**slack-archiver** is a simple tool to archive a Slack enterprise grid.

slack-archiver is built in [Go](https://golang.org/) and processes messages as a stream when possible to reduce memory resource requirements.

## Usage

Use `slack-archiver --version` to show the current version.

Below is the usage for the `slack-archiver` command.

```text
slack-archiver is a tool to archive a Slack enterprise grid.

Usage:
  slack-archiver [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  download    download data
  help        Help about any command
  list        list data
  version     show version

Flags:
  -h, --help   help for slack-archiver

Use "slack-archiver [command] --help" for more information about a command.
```

## Examples

Below is a trivial example.

```shell
$ bin/slack-archiver list teams --src "Hello\ World\ Slack\ export\ Jan\ 31\ 2020\ -\ Jul\ 29\ 2020.zip | jq -c .name"
{"name":"hello world"}
```

## Building

**slack-archiver** is written in pure Go, so the only dependency needed to compile the program is [Go](https://golang.org/).  Go can be downloaded from <https://golang.org/dl/>.

This project uses [direnv](https://direnv.net/) to manage environment variables and automatically adding the `bin` and `scripts` folder to the path.  Install direnv and hook it into your shell.  The use of `direnv` is optional as you can always call sc directly with `bin/sc`.

If using `macOS`, follow the `macOS` instructions below.

To build a binary for your local operating system you can use `make bin/slack-archiver`.  To build for a release, you can use `make build_release`.  Additionally, you can call `go build` directly to support specific use cases.

### macOS

You can install `go` on macOS using homebrew with `brew install go`.

To install `direnv` on `macOS` use `brew install direnv`.  If using bash, then add `eval \"$(direnv hook bash)\"` to the `~/.bash_profile` file .  If using zsh, then add `eval \"$(direnv hook zsh)\"` to the `~/.zshrc` file.

## Testing

**CLI**

To run CLI testes use `make test_cli`, which uses [shUnit2](https://github.com/kward/shunit2).  If you recive a `shunit2:FATAL Please declare TMPDIR with path on partition with exec permission.` error, you can modify the `TMPDIR` environment variable in line or with `export TMPDIR=<YOUR TEMP DIRECTORY HERE>`. For example:

```shell
TMPDIR="/usr/local/tmp" make test_cli
```

**Go**

To run Go tests use `make test_go` (or `bash scripts/test.sh`), which runs unit tests, `go vet`, `go vet with shadow`, [errcheck](https://github.com/kisielk/errcheck), [staticcheck](https://staticcheck.io/), and [misspell](https://github.com/client9/misspell).

## Contributing

We'd love to have your contributions!  Please see [CONTRIBUTING.md](CONTRIBUTING.md) for more info.

## Security

Please see [SECURITY.md](SECURITY.md) for more info.

## License

This project constitutes a work of the United States Government and is not subject to domestic copyright protection under 17 USC § 105.  However, because the project utilizes code licensed from contributors and other third parties, it therefore is licensed under the MIT License.  See LICENSE file for more information.
