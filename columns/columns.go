package columns

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Columns struct {
	Names           []string
	SymbolizedNames []string
}

// NamesString is a comma separated list of the column names.
func (c *Columns) NamesString() string {
	return strings.Join(c.Names, ", ")
}

// SymbolizedNamesString returns a comma separated list of
// the column names with a colon in front of each name.
func (c *Columns) SymbolizedNamesString() string {
	return strings.Join(c.SymbolizedNames, ", ")
}

// UpdatesString returns a comma separated list key = :value
// pairs, used for building named `UPDATE` statements.
func (c *Columns) UpdatesString() string {
	sets := []string{}
	for i := 0; i < len(c.Names); i++ {
		sets = append(sets, fmt.Sprintf("%s = %s", c.Names[i], c.SymbolizedNames[i]))
	}
	return strings.Join(sets, ", ")
}

// Add a column to the list.
func (c *Columns) Add(name string) {
	c.Names = append(c.Names, name)
	c.SymbolizedNames = append(c.SymbolizedNames, fmt.Sprintf(":%s", name))
}

// Remove a column from the list.
func (c *Columns) Remove(names ...string) {
	tmp := []string{}
	for _, name := range c.Names {
		found := false
		for _, n := range names {
			if n == name {
				found = true
				break
			}
		}
		if !found {
			tmp = append(tmp, name)
		}
	}
	c.Names = []string{}
	c.SymbolizedNames = []string{}
	for _, n := range tmp {
		c.Add(n)
	}
}

// ColumnsForStruct returns a Columns instance for
// the struct passed in.
func ColumnsForStruct(s interface{}) Columns {
	columns := Columns{}
	st := reflect.TypeOf(s)
	field_count := st.NumField()

	var w sync.WaitGroup
	w.Add(field_count)

	for i := 0; i < field_count; i++ {
		field := st.Field(i)
		go func(field reflect.StructField, columns *Columns, w *sync.WaitGroup) {

			defer w.Done()
			tag := field.Tag.Get("db")
			if tag == "" {
				tag = field.Name
			}
			if tag != "-" {
				columns.Add(tag)
			}
		}(field, &columns, &w)
	}

	w.Wait()
	return columns
}
