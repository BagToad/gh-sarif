# GH-SARIF

Interact with Code Scanning analysis and SARIF files. 

`gh-sarif` is a [GitHub CLI](https://github.com/cli/cli) extension. 

## Installation

```sh
gh extension install bagtoad/gh-sarif
```

## Usage

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

