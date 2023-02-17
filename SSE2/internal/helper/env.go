package helper

import "fmt"

func ConcatenateEnvVarPrefixes(prefixA string, prefixB string) string {
	return fmt.Sprintf("%s_%s", prefixA, prefixB)
}
