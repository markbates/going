package validators

import (
	"strings"

	"github.com/serenize/snaker"
)

func GenerateKey(s string) string {
	key := strings.Replace(s, " ", "", -1)
	key = strings.Replace(key, "-", "", -1)
	key = snaker.CamelToSnake(key)
	return key
}
