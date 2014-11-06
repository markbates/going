package columns

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"sync"
)

// Add a column to the list.
func (c *Columns) Add(names ...string) []*Column {
	ret := []*Column{}
	c.lock.Lock()
	for _, name := range names {
		xs := strings.Split(name, ",")
		col := c.Cols[xs[0]]
		if col == nil {
			col = &Column{
				Name:      xs[0],
				SelectSQL: xs[0],
				Readable:  true,
				Writeable: true,
			}

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
		ret = append(ret, col)
	}
	c.lock.Unlock()
	return ret
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
	SelectSQL string
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

func (c ReadableColumns) SelectString() string {
	xs := []string{}
	for _, t := range c.Cols {
		xs = append(xs, t.SelectSQL)
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}

type Columns struct {
	Cols map[string]*Column
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
		Cols: map[string]*Column{},
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
			cs := columns.Add(tag)
			c := cs[0]
			tag = field.Tag.Get("select")
			if tag != "" {
				log.Printf("tag: %s\n", tag)
				c.SelectSQL = tag
			}
		}
	}

	return columns
}
