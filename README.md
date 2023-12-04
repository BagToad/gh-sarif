# GH-SARIF

This [GitHub CLI](https://github.com/cli/cli) extension makes it easier to interact with GitHub's Code Scanning analyses endpoints, with some additional features added in.

## Features

- [COMPLETE] [`list` Code Scanning analyses](https://docs.github.com/en/rest/code-scanning/code-scanning?apiVersion=2022-11-28#list-code-scanning-analyses-for-a-repository)
- [WIP] [`view` Code Scanning analyses](https://docs.github.com/en/rest/code-scanning/code-scanning?apiVersion=2022-11-28#get-a-code-scanning-analysis-for-a-repository)
- [WIP] [`delete` Code Scanning analyses](https://docs.github.com/en/rest/code-scanning/code-scanning?apiVersion=2022-11-28#delete-a-code-scanning-analysis-from-a-repository).
- [WIP] SARIF format:
    - [WIP] `download` raw SARIF.
    - [WIP] `view` a processed SARIF summary (print results count, warnings, errors, and  ).

## Use Cases

- You need to script mass deleting Code Scanning analyses. 
- You find yourself regularly looking at Code Scanning analyses - perhaps you are integrating with Code Scanning. 
- You 