package columns_test

import (
	"testing"

	. "github.com/markbates/going/columns"
	"github.com/stretchr/testify/assert"
)

func TestColumnsForStruct(t *testing.T) {
	assert := assert.New(t)
	type foo struct {
		FirstName string `db:"first_name"`
		LastName  string
		Email     string `db:"email_address"`
		Unwanted  string `db:"-"`
	}

	f := foo{}
	columns := ColumnsForStruct(f)
	assert.Equal(columns.Names, []string{"first_name", "LastName", "email_address"})
	assert.Equal(columns.SymbolizedNames, []string{":first_name", ":LastName", ":email_address"})
	assert.Equal(columns.NamesString(), "first_name, LastName, email_address")
	assert.Equal(columns.SymbolizedNamesString(), ":first_name, :LastName, :email_address")
	assert.Equal(columns.UpdatesString(), "first_name = :first_name, LastName = :LastName, email_address = :email_address")

	columns.Add("bar")
	assert.Equal(columns.Names, []string{"first_name", "LastName", "email_address", "bar"})
	assert.Equal(columns.SymbolizedNames, []string{":first_name", ":LastName", ":email_address", ":bar"})
	assert.Equal(columns.NamesString(), "first_name, LastName, email_address, bar")
	assert.Equal(columns.SymbolizedNamesString(), ":first_name, :LastName, :email_address, :bar")
	assert.Equal(columns.UpdatesString(), "first_name = :first_name, LastName = :LastName, email_address = :email_address, bar = :bar")

	columns.Remove("bar", "email_address")
	assert.Equal(columns.Names, []string{"first_name", "LastName"})
	assert.Equal(columns.SymbolizedNames, []string{":first_name", ":LastName"})
	assert.Equal(columns.NamesString(), "first_name, LastName")
	assert.Equal(columns.SymbolizedNamesString(), ":first_name, :LastName")
	assert.Equal(columns.UpdatesString(), "first_name = :first_name, LastName = :LastName")
}
