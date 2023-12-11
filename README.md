# GH-SARIF

`gh-sarif` is a [GitHub CLI](https://github.com/cli/cli) extension to interact with Code Scanning analyses.

This is a work in progress (and mostly just for fun)! Anything may change or be removed at any time. 

## Features

- [COMPLETE] [`list` Code Scanning analyses](https://docs.github.com/en/rest/code-scanning/code-scanning?apiVersion=2022-11-28#list-code-scanning-analyses-for-a-repository)
- [WIP] [`view` Code Scanning analyses](https://docs.github.com/en/rest/code-scanning/code-scanning?apiVersion=2022-11-28#get-a-code-scanning-analysis-for-a-repository)
- [WIP] [`delete` Code Scanning analyses](https://docs.github.com/en/rest/code-scanning/code-scanning?apiVersion=2022-11-28#delete-a-code-scanning-analysis-from-a-repository).
- [WIP] SARIF format:
    - [WIP] `download` raw SARIF.
    - [WIP] `view` a processed SARIF summary (print results count, warnings, errors, and  ).

## Use Cases

There are a few use cases I'll hope to fill with this:

- You find yourself regularly looking at Code Scanning analyses - perhaps you are integrating with Code Scanning. 
- You need to delete all Code Scanning analyses in a set. 
- You have an automation in GitHub Actions (or somewhere else with the GitHub CLI installed) that requires you to work with Code Scanning analyses. 