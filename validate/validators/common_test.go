package validators_test

import (
	"testing"

	"github.com/markbates/going/validate/validators"
	"github.com/stretchr/testify/require"
)

func Test_GenerateKey(t *testing.T) {
	r := require.New(t)

	r.Equal("foo", validators.GenerateKey("Foo"))
	r.Equal("created_at", validators.GenerateKey("CreatedAt"))
	r.Equal("created_at", validators.GenerateKey("Created At"))
	r.Equal("person_id", validators.GenerateKey("PersonID"))
	r.Equal("content_type", validators.GenerateKey("Content-Type"))
}
