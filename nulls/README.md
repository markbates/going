# github.com/markbates/going/nulls

This package should be used in place of the built-in null types in the `sql` package.

The real benefit of this packages comes in its implementation of `MarshalJSON` and `UnmarshalJSON` to properly encode/decode `null` values.

## Installation

``` bash
$ go get github.com/markbates/going/nulls
```

## Supported Datatypes

* `string` (`nulls.NullString`) - Replaces `sql.NullString`
* `int64` (`nulls.NullInt64`) - Replaces `sql.NullInt64`
* `float64` (`nulls.NullFloat64`) - Replaces `sql.NullFloat64`
* `bool` (`nulls.NullBool`) - Replaces `sql.NullBool`

Additionally this package provides the following extra data types:

* `time.Time` (`nulls.NullTime`)
