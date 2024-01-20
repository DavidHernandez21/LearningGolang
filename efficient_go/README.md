
## check the sizes of your CPU caches
```shell
sudo dmidecode -t cache
```

## Display information about the CPU architecture
```shell
lscpu
```


## go build flags

```
There are many different configuration options for the Go build process. The first
large batch of options can be passed through go build -ldflags="<flags>", which
represents linker command options (the ld prefix traditionally stands for Linux
linker). For example:
	• We can omit the DWARF table, thus reducing the binary size using
	-ldflags="-w" (recommended for production build if you don’t use debuggers
	there).
	• We can further reduce the size with -ldflags= "-s -w", removing the DWARF
	and symbols tables with other debug information. I would not recommend the
	latter option, as non-DWARF elements allow important runtime routines, like
	gathering profiles.
Similarly, go build -gcflags="<flags>" represents Go compiler options (gc stands
for Go Compiler; don’t confuse it with GC, which means garbage collection, as
explained in “Garbage Collection” on page 185). For example:
• -gcflags="-S" prints Go Assembly from the source code.
• -gcflags="-N" disables all compiler optimizations.
```

## go build/run with heap allocation info
```shell
 go run -gcflags='-m' .\main.go
 full debug
 go run -gcflags='-m -m' .\main.go
```

## go build/run bound check info

```shell
 go run -gcflags="-d=ssa/check_bce" .\main.go
```

## point the pprof tool directly to the profiler URL to avoid the manual process of downloading the file

```shell
 go tool pprof -http :8080 http://<address>/debug/pprof/heap
```

> If you have multiple allocations in a single function, it is often useful to analyze the
	heap profile in lines granularity (add the &g=lines URL parameter in the web
	viewer).

## Comparing and Aggregating Profiles

```shell
 go tool pprof heap-AB.pprof -base heap-B.pprof
 go tool pprof heap-AB.pprof -base heap-B.pprof -diff_base
```
