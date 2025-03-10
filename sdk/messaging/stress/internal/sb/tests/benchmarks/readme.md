## Generating

`go test -memprofile mem.out -bench .`

## Visualizing

Run:
* `sudo apt install graphviz`
* `go tool pprof -http localhost:8000 -base mem.out.before_fix mem.out.after_fix`
