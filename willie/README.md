# github.com/markbates/going/willie

A helper for cleaning up HTTP tests.

## Installation

```bash
$ go get https://github.com/markbates/going/willie
```

## Usage

```go
func Test_Get(t *testing.T) {
	a := require.New(t)

	w := willie.New(app())

	res := w.Get("/get", nil)
}
```
