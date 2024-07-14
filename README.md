# GH-SARIF

Interact with Code Scanning analysis and SARIF files. 

`gh-sarif` is a [GitHub CLI](https://github.com/cli/cli) extension. 

## Installation

```sh
gh extension install bagtoad/gh-sarif
```

## Usage

```
Usage:
  gh sarif [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a code scanning analysis
  help        Help about any command
  list        List analyses for a repository
  upload      Upload a SARIF file to a repo
  view        View analysis results

Flags:
  -h, --help          help for gh-sarif
  -j, --json          Output JSON instead of text (includes additional fields)
  -R, --repo string   GitHub repository (format: owner/repo)

Use "gh sarif [command] --help" for more information about a command.
```

### List Analyses for a Repository

```sh
gh sarif list
```
### View Analysis Results in a Table

```sh
gh sarif view <analysis-id>
```

### View Analysis Results as SARIF

```sh
gh sarif view <analysis-id> --sarif
```

### View Analysis Results as CSV

```sh
gh sarif view <analysis-id> --csv
```

### View Analysis Results from a Local SARIF File

```sh
gh sarif view <path-to-sarif-file>
```

### Upload a SARIF File to GitHub Code Scanning

```sh
gh sarif upload <commit-sha> <ref> <path-to-sarif-file>
```

### Delete an Analysis

```sh
gh sarif delete <analysis-id>
```

### Delete Multiple Analyses

```sh
gh sarif delete <analysis-id> <analysis-id> <analysis-id>...
```

### Delete All Analyses in the set Except the Last

```sh
gh sarif delete <analysis-id> --delete-all
```

### Delete All Analyses in the set, Including the Last

```sh
gh sarif delete <analysis-id> --delete-all --confirm-delete
```

or 

```sh
gh sarif delete <analysis-id> --purge
```

