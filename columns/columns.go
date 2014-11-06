package columns

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

// Add a column to the list.
func (c *Columns) Add(names ...string) {
	c.lock.Lock()
	for _, name := range names {
		xs := strings.Split(name, ",")
		col := Column{}
		col.Name = xs[0]
		if c.Cols[col.Name].Name == "" {
			col.Readable = true
			col.Writeable = true

			if len(xs) > 1 {
				if xs[1] == "*readonly" {
					col.Writeable = false
				}
				if xs[1] == "*writeonly" {
					col.Readable = false
				}
			}

			c.Cols[col.Name] = col
		}
	}
	c.lock.Unlock()
}

// Remove a column from the list.
func (c *Columns) Remove(names ...string) {
	for _, name := range names {
		xs := strings.Split(name, ",")
		name = xs[0]
		delete(c.Cols, name)
	}
}

type Column struct {
	Name      string
	Writeable bool
	Readable  bool
}

func (c Column) UpdateString() string {
	return fmt.Sprintf("%s = :%s", c.Name, c.Name)
}

type WriteableColumns struct {
	Columns
}

func (c WriteableColumns) UpdateString() string {
	xs := []string{}
	for _, t := range c.Cols {
		xs = append(xs, t.UpdateString())
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}

type ReadableColumns struct {
	Columns
}

type Columns struct {
	Cols map[string]Column
	lock *sync.RWMutex
}

func (c Columns) Writeable() *WriteableColumns {
	w := &WriteableColumns{NewColumns()}
	for _, col := range c.Cols {
		if col.Writeable {
			w.Cols[col.Name] = col
		}
	}
	return w
}

func (c Columns) Readable() *ReadableColumns {
	w := &ReadableColumns{NewColumns()}
	for _, col := range c.Cols {
		if col.Readable {
			w.Cols[col.Name] = col
		}
	}
	return w
}

func (c Columns) String() string {
	xs := []string{}
	for _, t := range c.Cols {
		xs = append(xs, t.Name)
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}

func (c Columns) SymbolizedString() string {
	xs := []string{}
	for _, t := range c.Cols {
		xs = append(xs, ":"+t.Name)
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}

func NewColumns() Columns {
	return Columns{
		lock: &sync.RWMutex{},
		Cols: map[string]Column{},
	}
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
