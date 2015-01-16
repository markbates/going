package nulls_test

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	. "github.com/markbates/going/nulls"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

type Foo struct {
	ID         NullInt64     `json:"id" db:"id"`
	Name       NullString    `json:"name" db:"name"`
	Alive      NullBool      `json:"alive" db:"alive"`
	Price      NullFloat64   `json:"price" db:"price"`
	Birth      NullTime      `json:"birth" db:"birth"`
	Price32    NullFloat32   `json:"price32" db:"price32"`
	Bytes      NullByteSlice `json:"bytes" db:"bytes"`
	IntType    NullInt       `json:"intType" db:"int_type"`
	Int32Type  NullInt32     `json:"int32Type" db:"int32_type"`
	UInt32Type NullUInt32    `json:"uint32Type" db:"uint32_type"`
}

const schema = `CREATE TABLE "main"."foos" (
	 "id" integer,
	 "name" text,
	 "alive" integer,
	 "price" float,
	 "birth" timestamp,
	 "price32" float,
	 "bytes"  blob,
	 "int_type" integer,
	 "int32_type" integer,
	 "uint32_type" integer
);`

var now = time.Now()

func newValidFoo() Foo {
	return Foo{
		ID:         NewNullInt64(1),
		Name:       NewNullString("Mark"),
		Alive:      NewNullBool(true),
		Price:      NewNullFloat64(9.99),
		Birth:      NewNullTime(now),
		Price32:    NewNullFloat32(3.33),
		Bytes:      NewNullByteSlice([]byte("Byte Slice")),
		IntType:    NewNullInt(2),
		Int32Type:  NewNullInt32(3),
		UInt32Type: NewNullUInt32(5),
	}
}

func TestNullTypesMarshalProperly(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	f := newValidFoo()

	ti, _ := json.Marshal(now)
	ba, _ := json.Marshal(f.Bytes)
	jsonString := fmt.Sprintf(`{"id":1,"name":"Mark","alive":true,"price":9.99,"birth":%s,"price32":3.33,"bytes":%s,"intType":2,"int32Type":3,"uint32Type":5}`, ti, ba)

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
	assert.Equal(f.Bytes.ByteSlice, ba)
	assert.Equal(f.IntType.Int, 2)
	assert.Equal(f.Int32Type.Int32, 3)
	assert.Equal(f.UInt32Type.UInt32, uint32(5))

	// check marshalling nulls works:
	f = Foo{}
	jsonString = `{"id":null,"name":null,"alive":null,"price":null,"birth":null,"price32":null,"bytes":null,"intType":null,"int32Type":null,"uint32Type":null}`
	data, _ = json.Marshal(f)
	assert.Equal(string(data), jsonString)

	f = Foo{}
	json.NewDecoder(strings.NewReader(jsonString)).Decode(&f)
	assert.Equal(f.ID.Int64, 0)
	assert.False(f.ID.Valid)
	assert.Equal(f.Name.String, "")
	assert.False(f.Name.Valid)
	assert.Equal(f.Alive.Bool, false)
	assert.False(f.Alive.Valid)
	assert.Equal(f.Price.Float64, 0)
	assert.False(f.Price.Valid)
	assert.Equal(f.Birth.Time, time.Time{})
	assert.False(f.Birth.Valid)
	assert.Equal(f.Price32.Float32, 0)
	assert.False(f.Price32.Valid)
	assert.Equal(f.Bytes.ByteSlice, []byte(nil))
	assert.False(f.Bytes.Valid)
	assert.Equal(f.IntType.Int, 0)
	assert.False(f.IntType.Valid)
	assert.Equal(f.Int32Type.Int32, 0)
	assert.False(f.Int32Type.Valid)
	assert.Equal(f.UInt32Type.UInt32, uint32(0))
	assert.False(f.UInt32Type.Valid)
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
		assert.False(f.Bytes.Valid)
		assert.False(f.IntType.Valid)
		assert.False(f.Int32Type.Valid)
		assert.False(f.UInt32Type.Valid)
		assert.Equal(f.Birth.Time.UnixNano(), time.Time{}.UnixNano())
		assert.Equal(f.ID.Int64, 0)
		assert.Equal(f.Name.String, "")
		assert.Equal(f.Price.Float64, 0)
		assert.Equal(f.Price32.Float32, 0)
		assert.Equal(f.Bytes.ByteSlice, []byte(nil))
		assert.Equal(f.IntType.Int, 0)
		assert.Equal(f.Int32Type.Int32, 0)
		assert.Equal(f.UInt32Type.UInt32, uint32(0))
		tx.Rollback()

		tx, err = db.Beginx()
		assert.NoError(err)

		f = newValidFoo()
		tx.NamedExec("INSERT INTO foos (id, name, alive, price, birth, price32, bytes, int_type, int32_type, uint32_type) VALUES (:id, :name, :alive, :price, :birth, :price32, :bytes, :int_type, :int32_type, :uint32_type)", &f)
		f = Foo{}
		tx.Get(&f, "select * from foos")
		assert.True(f.Alive.Valid)
		assert.True(f.Birth.Valid)
		assert.True(f.ID.Valid)
		assert.True(f.Name.Valid)
		assert.True(f.Price.Valid)
		assert.True(f.Alive.Bool)
		assert.True(f.Price32.Valid)
		assert.True(f.Bytes.Valid)
		assert.True(f.IntType.Valid)
		assert.True(f.Int32Type.Valid)
		assert.True(f.UInt32Type.Valid)
		assert.Equal(f.Birth.Time.UnixNano(), now.UnixNano())
		assert.Equal(f.ID.Int64, 1)
		assert.Equal(f.Name.String, "Mark")
		assert.Equal(f.Price.Float64, 9.99)
		assert.Equal(f.Price32.Float32, 3.33)
		assert.Equal(f.Bytes.ByteSlice, []byte("Byte Slice"))
		assert.Equal(f.IntType.Int, 2)
		assert.Equal(f.Int32Type.Int32, 3)
		assert.Equal(f.UInt32Type.UInt32, uint32(5))

		tx.Rollback()
	})
}
