# LearningGolang
In this repo I will put code that I have wrote following tutorials and reading blogs---all about golang

## golangci-lint
[docs](https://golangci-lint.run/usage/configuration)

- run the golangci-lint container
    + docker run --rm -it -v $PWD:/app -w /app golangci/golangci-lint bash
- To see a list of enabled by your configuration linters
    + golangci-lint linters
- Run the linters
    + golangci-lint -v run // verbose mode

## gitleaks (maybe better using it with pre-commit)
[docs](https://github.com/zricethezav/gitleaks)

- docker pull zricethezav/gitleaks:latest
- docker run --rm -it -v $PWD:/path  zricethezav/gitleaks:latest detect --source=path
