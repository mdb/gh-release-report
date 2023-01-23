[![CI/CD](https://github.com/mdb/gh-release-report/actions/workflows/cicd.yaml/badge.svg)](https://github.com/mdb/gh-release-report/actions/workflows/cicd.yaml)

# gh-release-report

A `gh` CLI extension that reports and visualizes a GitHub release's total
download count, as well as the individual download counts for each of its assets.

![demo](demo.gif)

## Usage

Inspect the download counts associated with a particular GitHub release:

```
gh release-report \
  --repo "cli/cli" \
  --tag "v1.0.0"
```

By default, the latest release associated with the current working directory's GitHub repository is targeted.

## Installation

Install the `gh` CLI [for your platform](https://github.com/cli/cli#installation). For example, on Mac OS:

```
brew install gh
```

Install the latest `release-report` extension from [its releases](https://github.com/mdb/gh-release-report/releases):

```
gh extension install mdb/gh-release-report
```

## Development

Build and test `gh-release-report` locally:

```
make
```

Install a locally built `gh-release-report` for use as `gh release-report`:

```
make install
```
