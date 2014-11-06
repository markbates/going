package columns_test

import (
	"testing"

	. "github.com/markbates/going/columns"
	"github.com/stretchr/testify/assert"
)

type foo struct {
	FirstName string `db:"first_name" select:"first_name as f"`
	LastName  string
	Unwanted  string `db:"-"`
	ReadOnly  string `db:"read,*readonly"`
	WriteOnly string `db:"write,*writeonly"`
}

func Test_Column_UpdateString(t *testing.T) {
	a := assert.New(t)
	c := Column{Name: "foo"}
	a.Equal(c.UpdateString(), "foo = :foo")
}

func Test_Columns_UpdateString(t *testing.T) {
	a := assert.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := ColumnsForStruct(f)
		u := c.Writeable().UpdateString()
		a.Equal(u, "LastName = :LastName, first_name = :first_name, write = :write")
	}
}

func Test_Columns_WriteableString(t *testing.T) {
	a := assert.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := ColumnsForStruct(f)
		u := c.Writeable().String()
		a.Equal(u, "LastName, first_name, write")
	}
}

func Test_Columns_ReadableString(t *testing.T) {
	a := assert.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := ColumnsForStruct(f)
		u := c.Readable().String()
		a.Equal(u, "LastName, first_name, read")
	}
}

func Test_Columns_Readable_SelectString(t *testing.T) {
	a := assert.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := ColumnsForStruct(f)
		u := c.Readable().SelectString()
		a.Equal(u, "LastName, first_name as f, read")
	}
}

func Test_Columns_WriteableString_Symbolized(t *testing.T) {
	a := assert.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := ColumnsForStruct(f)
		u := c.Writeable().SymbolizedString()
		a.Equal(u, ":LastName, :first_name, :write")
	}
}

func Test_Columns_ReadableString_Symbolized(t *testing.T) {
	a := assert.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := ColumnsForStruct(f)
		u := c.Readable().SymbolizedString()
		a.Equal(u, ":LastName, :first_name, :read")
	}
}
func Test_Columns_Basics(t *testing.T) {
	a := assert.New(t)

	for _, f := range []interface{}{foo{}, &foo{}} {
		c := ColumnsForStruct(f)
		a.Equal(len(c.Cols), 4)
		a.Equal(c.Cols["first_name"], &Column{Name: "first_name", Writeable: true, Readable: true, SelectSQL: "first_name as f"})
		a.Equal(c.Cols["LastName"], &Column{Name: "LastName", Writeable: true, Readable: true, SelectSQL: "LastName"})
		a.Equal(c.Cols["read"], &Column{Name: "read", Writeable: false, Readable: true, SelectSQL: "read"})
		a.Equal(c.Cols["write"], &Column{Name: "write", Writeable: true, Readable: false, SelectSQL: "write"})
	}
}

func Test_Columns_Add(t *testing.T) {
	a := assert.New(t)

	for _, f := range []interface{}{foo{}, &foo{}} {
		c := ColumnsForStruct(f)
		a.Equal(len(c.Cols), 4)
		c.Add("foo", "first_name")
		a.Equal(len(c.Cols), 5)
		a.Equal(c.Cols["foo"], &Column{Name: "foo", Writeable: true, Readable: true, SelectSQL: "foo"})
	}
}

func Test_Columns_Remove(t *testing.T) {
	a := assert.New(t)

	for _, f := range []interface{}{foo{}, &foo{}} {
		c := ColumnsForStruct(f)
		a.Equal(len(c.Cols), 4)
		c.Remove("foo", "first_name")
		a.Equal(len(c.Cols), 3)
	}
}

// func Test_ColumnsForStruct(t *testing.T) {
// 	t.Parallel()
// 	assert := assert.New(t)
//
// 	f := foo{}
// 	columns := ColumnsForStruct(f)
// 	assert.Equal(columns.Names, []string{"first_name", "LastName", "email_address", "read"})
// 	assert.Equal(columns.SymbolizedNames, []string{":first_name", ":LastName", ":email_address", ":read"})
// 	assert.Equal(columns.NamesString(), "first_name, LastName, email_address, read")
// 	assert.Equal(columns.SymbolizedNamesString(), ":first_name, :LastName, :email_address, :read")
// 	assert.Equal(columns.UpdatesString(), "first_name = :first_name, LastName = :LastName, email_address = :email_address")
//
// 	columns.Add("bar")
// 	assert.Equal(columns.Names, []string{"first_name", "LastName", "email_address", "read", "bar"})
// 	assert.Equal(columns.SymbolizedNames, []string{":first_name", ":LastName", ":email_address", "read", ":bar"})
// 	assert.Equal(columns.NamesString(), "first_name, LastName, email_address, read, bar")
// 	assert.Equal(columns.SymbolizedNamesString(), ":first_name, :LastName, :email_address, :read, :bar")
// 	assert.Equal(columns.UpdatesString(), "first_name = :first_name, LastName = :LastName, email_address = :email_address, bar = :bar")
//
// 	columns.Remove("bar", "email_address")
// 	assert.Equal(columns.Names, []string{"first_name", "LastName", "read"})
// 	assert.Equal(columns.SymbolizedNames, []string{":first_name", ":LastName", "read"})
// 	assert.Equal(columns.NamesString(), "first_name, LastName, read")
// 	assert.Equal(columns.SymbolizedNamesString(), ":first_name, :LastName, :read")
// 	assert.Equal(columns.UpdatesString(), "first_name = :first_name, LastName = :LastName")
// }
//
// func Test_Columns_Add_Duplicates(t *testing.T) {
// 	t.Parallel()
// 	a := assert.New(t)
//
// 	columns := NewColumns()
// 	columns.Add("foo")
// 	a.Equal(columns.Names, []string{"foo"})
//
// 	// adding the same column again should have no effect:
// 	columns.Add("foo")
// 	a.Equal(columns.Names, []string{"foo"})
// }
//
// func Test_ColumnsForStruct_WithPointer(t *testing.T) {
// 	t.Parallel()
// 	assert := assert.New(t)
// 	f := &foo{}
//
// 	columns := ColumnsForStruct(f)
// 	assert.Equal(columns.Names, []string{"first_name", "LastName", "email_address"})
// 	assert.Equal(columns.SymbolizedNames, []string{":first_name", ":LastName", ":email_address"})
// }
