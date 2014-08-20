package nulls_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/markbates/deano/nulls"
	"github.com/stretchr/testify/assert"
)

type foo struct {
	ID    NullInt64   `json:"id" db:"id"`
	Name  NullString  `json:"name" db:"name"`
	Alive NullBool    `json:"alive" db:"alive"`
	Price NullFloat64 `json:"price" db:"price"`
	Birth NullTime    `json:"birth" db:"birth"`
}

func TestNullTypesMarshalProperly(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	f := foo{
		ID:    NewNullInt64(1),
		Name:  NewNullString("Mark"),
		Alive: NewNullBool(true),
		Price: NewNullFloat64(9.99),
		Birth: NewNullTime(now),
	}

	ti, _ := json.Marshal(now)
	jsonString := fmt.Sprintf(`{"id":1,"name":"Mark","alive":true,"price":9.99,"birth":%s}`, ti)

	// check marshalling to json works:
	data, _ := json.Marshal(f)
	assert.Equal(string(data), jsonString)

	// check unmarshalling from json works:
	f = foo{}
	json.NewDecoder(strings.NewReader(jsonString)).Decode(&f)
	assert.Equal(f.ID.Int64, 1)
	assert.Equal(f.Name.String, "Mark")
	assert.Equal(f.Alive.Bool, true)
	assert.Equal(f.Price.Float64, 9.99)
	assert.Equal(f.Birth.Time, now)

	// check marshalling nulls works:
	f = foo{}
	jsonString = `{"id":null,"name":null,"alive":false,"price":null,"birth":null}`
	data, _ = json.Marshal(f)
	assert.Equal(string(data), jsonString)

	f = foo{}
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
}

func TestNullTypeSaveAndRetrieveProperly(t *testing.T) {
	t.Skip("Need to fill this in with proper tests talking to a DB")
}
