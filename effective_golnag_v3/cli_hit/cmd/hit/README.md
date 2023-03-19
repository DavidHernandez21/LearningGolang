## go build strip debugging information
go build -ldflags "-s -w" -o hit_s_w.exe

## extra compression usin UPX
upx --best -o hit_upx.exe .\hit_s_w.exe

## test the flags.go using the tag //go:build cli
go test -v -run=".*" -tags="cli"
