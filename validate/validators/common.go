package validators

import "strings"

func generateKey(s string) string {
	key := strings.ToLower(s)
	key = strings.Replace(key, " ", "_", -1)
	return key
}
