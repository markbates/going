# github.com/markbates/going/emu

A collection of middleware for [echo](https://github.com/labstack/echo)

## Installation

```bash
$ go get https://github.com/markbates/going/emu
```

## Usage

```go
func Test_Get(t *testing.T) {
	a := require.New(t)

	w := willie.New(app())

	res := w.Get("/get", nil)
}
```
