go build .\main.g

// --workers defaults to runtime.NumCPU()
.\main.exe --workers=4

go tool trace .\trace.out
