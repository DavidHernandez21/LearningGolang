## go build/run with heap allocation info
```shell
    go run -gcflags='-m' .\main.go
    # full debug
    go run -gcflags='-m -m' .\main.go
```

## go build/run bound check info

```shell
    go run -gcflags="-d=ssa/check_bce" .\main.go
```

## get a 30 seconds profile

```shell
    # run the program
    go run .\main.go
    # open another terminal and run the following command
    curl -o default.pgo http://localhost:6060/debug/pprof/profile?seconds=30
    # open the default.pgo with go tool pprof
    go tool pprof default.pgo

```

## get 10 seconds trace

```shell
    # run the program
    go run .\main.go
    # open another terminal and run the following command
    curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=10
    # open the trace.out with go tool trace
    go tool trace trace.out
```
