# github.com/markbates/going/columns

This package helps you get a list of database columns for a struct, and then manage those columns. This is useful for building SQL statements that require the names of columns, such as `SELECT`, `INSERT`, or `UPDATE` statements.

## Installation

``` bash
$ go get github.com/markbates/going/columns
```

## Usage

```go
package main

import (
	. "github.com/markbates/going/columns"
)

type User struct {
	FirstName string `db:"first_name"`
	LastName  string
	Email     string `db:"email_address"`
	Unwanted  string `db:"-"`
}

func main() {
	u := User{}
	columns := ColumnsForStruct(u)
	columns.Names //[]string{"first_name", "LastName", "email_address"}
	columns.SymbolizedNames // []string{":first_name", ":LastName", ":email_address"}
	columns.NamesString() // "first_name, LastName, email_address"
	columns.SymbolizedNamesString() // ":first_name, :LastName, :email_address"
	columns.UpdatesString() // "first_name = :first_name, LastName = :LastName, email_address = :email_address"

	columns.Add("bar")
	columns.Names // []string{"first_name", "LastName", "email_address", "bar"}
	columns.SymbolizedNames // []string{":first_name", ":LastName", ":email_address", ":bar"}
	columns.NamesString() // "first_name, LastName, email_address, bar"
	columns.SymbolizedNamesString() // ":first_name, :LastName, :email_address, :bar"
	columns.UpdatesString() // "first_name = :first_name, LastName = :LastName, email_address = :email_address, bar = :bar"

	columns.Remove("bar", "email_address")
	columns.Names // []string{"first_name", "LastName"}
	columns.SymbolizedNames // []string{":first_name", ":LastName"}
	columns.NamesString() // "first_name, LastName"
	columns.SymbolizedNamesString() // ":first_name, :LastName"
	columns.UpdatesString() // "first_name = :first_name, LastName = :LastName"
}
```
