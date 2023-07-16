## Run golangci-lint
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint golangci-lint -v run
