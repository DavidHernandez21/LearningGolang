## Test coverage
go test -coverprofile cover.out
go tool cover -html=cover.out
### avoid the browser
go tool cover -html=cover.out -o cover.html 