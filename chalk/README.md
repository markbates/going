# github.com/markbates/going/chalk

Quickly and easily create a pool of `go routines`.

## Installation

```bash
$ go get https://github.com/markbates/going/chalk
```

## Usage

```go
// create a pool of 10 go routines
tasks, errs := chalk.New(10)

tasks <- func() error {
  // do some work here
  return nil
}

// handle err
tasks <- func() error {
  return errors.New("boom!")
}

err := <-errs
```
