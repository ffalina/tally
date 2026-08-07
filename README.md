# tally

A tiny math expression solver written in Go. Parses and evaluates expressions like `3 + 4 * 2 / (1 - 5)^2`, with step-by-step output.

Supports:
- Arithmetic: `+ - * / ^`, parentheses, factorial (`5!`)
- Functions: `sqrt`, `abs`, `floor`, `ceil`, `round`, `log` (natural, or `log(x, base)`), `log10`
- Trig (degrees): `sin`, `cos`, `tan`, `asin`, `acos`, `atan`
- Combinatorics: `C(n, r)`, `P(n, r)`
- Constants: `pi`, `e`
- Linear equations with one variable (`x`), e.g. `2x + 3 = 11`

## Usage

```sh
go run ./src
```

## Test

```sh
go test ./src
```

## License

MIT
