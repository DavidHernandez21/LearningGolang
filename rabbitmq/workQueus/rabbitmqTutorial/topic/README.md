# Topic Exchange
- `*` (star) can substitute for exactly one word.
- `#` (star) can substitute for exactly one word.

## Consumer1
`go run .\main.go --routingKey=app.* --routingKey=db.critical.#`

## Consumer2
`go run .\main.go --routingKey=db.*`

## Producer
`go run .\main.go --mex='db.error'`
`go run .\main.go --mex='app.warning' --routingKey='app.warning'`
