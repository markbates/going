package nulls_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	. "github.com/markbates/going/nulls"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type Foo struct {
	ID      NullInt64   `json:"id" db:"id"`
	Name    NullString  `json:"name" db:"name"`
	Alive   NullBool    `json:"alive" db:"alive"`
	Price   NullFloat64 `json:"price" db:"price"`
	Birth   NullTime    `json:"birth" db:"birth"`
	Price32 NullFloat32 `json:"price32" db:"price32"`
}

const schema = `CREATE TABLE "main"."foos" (
	 "id" integer,
	 "name" text,
	 "alive" integer,
	 "price" float,
	 "birth" timestamp,
	 "price32" float
);`

var now = time.Now()

func newValidFoo() Foo {
	return Foo{
		ID:      NewNullInt64(1),
		Name:    NewNullString("Mark"),
		Alive:   NewNullBool(true),
		Price:   NewNullFloat64(9.99),
		Birth:   NewNullTime(now),
		Price32: NewNullFloat32(3.33),
	}
}

func TestNullTypesMarshalProperly(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	f := newValidFoo()

	ti, _ := json.Marshal(now)
	jsonString := fmt.Sprintf(`{"id":1,"name":"Mark","alive":true,"price":9.99,"birth":%s,"price32":3.33}`, ti)

	// check marshalling to json works:
	data, _ := json.Marshal(f)
	assert.Equal(string(data), jsonString)

	// check unmarshalling from json works:
	f = Foo{}
	json.NewDecoder(strings.NewReader(jsonString)).Decode(&f)
	assert.Equal(f.ID.Int64, 1)
	assert.Equal(f.Name.String, "Mark")
	assert.Equal(f.Alive.Bool, true)
	assert.Equal(f.Price.Float64, 9.99)
	assert.Equal(f.Birth.Time, now)
	assert.Equal(f.Price32.Float32, 3.33)

	// check marshalling nulls works:
	f = Foo{}
	jsonString = `{"id":null,"name":null,"alive":false,"price":null,"birth":null,"price32":null}`
	//jsonString = `{"id":null,"name":null,"alive":false,"price":null,"birth":null}`
	data, _ = json.Marshal(f)
	assert.Equal(string(data), jsonString)

	f = Foo{}
	json.NewDecoder(strings.NewReader(jsonString)).Decode(&f)
	assert.Equal(f.ID.Int64, 0)
	assert.False(f.ID.Valid)
	assert.Equal(f.Name.String, "")
	assert.False(f.Name.Valid)
	assert.False(f.Alive.Bool)
	assert.True(f.Alive.Valid)
	assert.Equal(f.Price.Float64, 0)
	assert.False(f.Price.Valid)
	assert.Equal(f.Birth.Time, time.Time{})
	assert.False(f.Birth.Valid)
	assert.Equal(f.Price32.Float32, 0)
	assert.False(f.Price32.Valid)
}

func initDB(f func(db *sqlx.DB)) {
	os.Remove("./foo.db")
	db, _ := sqlx.Open("sqlite3", "./foo.db")
	db.MustExec(schema)
	f(db)
	os.Remove("./foo.db")
}
func TestNullTypeSaveAndRetrieveProperly(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	initDB(func(db *sqlx.DB) {
		tx, err := db.Beginx()
		assert.NoError(err)
		tx.Exec("insert into foos")

		f := Foo{}
		tx.Get(&f, "select * from foos")
		assert.False(f.Alive.Valid)
		assert.False(f.Birth.Valid)
		assert.False(f.ID.Valid)
		assert.False(f.Name.Valid)
		assert.False(f.Price.Valid)
		assert.False(f.Alive.Bool)
		assert.False(f.Price32.Valid)
		assert.Equal(f.Birth.Time.UnixNano(), time.Time{}.UnixNano())
		assert.Equal(f.ID.Int64, 0)
		assert.Equal(f.Name.String, "")
		assert.Equal(f.Price.Float64, 0)
		assert.Equal(f.Price32.Float32, 0)
		tx.Rollback()

		tx, err = db.Beginx()
		assert.NoError(err)

		f = newValidFoo()
		tx.NamedExec("INSERT INTO foos (id, name, alive, price, birth, price32) VALUES (:id, :name, :alive, :price, :birth, :price32)", &f)
		f = Foo{}
		tx.Get(&f, "select * from foos")
		assert.True(f.Alive.Valid)
		assert.True(f.Birth.Valid)
		assert.True(f.ID.Valid)
		assert.True(f.Name.Valid)
		assert.True(f.Price.Valid)
		assert.True(f.Alive.Bool)
		assert.True(f.Price32.Valid)
		assert.Equal(f.Birth.Time.UnixNano(), now.UnixNano())
		assert.Equal(f.ID.Int64, 1)
		assert.Equal(f.Name.String, "Mark")
		assert.Equal(f.Price.Float64, 9.99)
		assert.Equal(f.Price32.Float32, 3.33)

		tx.Rollback()
	})
}
