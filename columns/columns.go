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
	keys            map[string]string
	lock            *sync.RWMutex
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
func (c *Columns) Add(names ...string) {
	c.lock.Lock()
	for _, name := range names {
		if c.keys[name] == "" {
			c.Names = append(c.Names, name)
			c.SymbolizedNames = append(c.SymbolizedNames, fmt.Sprintf(":%s", name))
			c.keys[name] = name
		}
	}
	c.lock.Unlock()
}

// Remove a column from the list.
func (c *Columns) Remove(names ...string) {
	tmp := []string{}
	for _, name := range c.Names {
		delete(c.keys, name)
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

func NewColumns() Columns {
	return Columns{keys: map[string]string{}, lock: &sync.RWMutex{}}
}

// ColumnsForStruct returns a Columns instance for
// the struct passed in.
func ColumnsForStruct(s interface{}) Columns {
	columns := NewColumns()
	st := reflect.TypeOf(s)
	if st.Kind().String() == "ptr" {
		st = reflect.ValueOf(s).Elem().Type()
	}
	field_count := st.NumField()

	for i := 0; i < field_count; i++ {
		field := st.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			tag = field.Name
		}
		if tag != "-" {
			columns.Add(tag)
		}
	}

	return columns
}
