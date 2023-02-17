package helper

import "fmt"

// ConcatenateEnvVarPrefixes склеивает два префикса названий переменных
// окружения.
func ConcatenateEnvVarPrefixes(prefixA string, prefixB string) string {
	return fmt.Sprintf("%s_%s", prefixA, prefixB)
}
