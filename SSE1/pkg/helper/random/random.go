package random

import (
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
	"github.com/vault-thirteen/auxie/random"
)

// Generates a unique random Marker using a mixed Technique including:
//   - An UUID of the fourth Version,
//   - Some random Bytes received from the random Numbers' Generator of the O.S.
func CreateUniqueMarker() (marker string, err error) {
	const (
		Separator = "-"
		Void      = ""
	)
	var uid uuid.UUID
	uid, err = uuid.NewRandom()
	if err != nil {
		return
	}
	marker = uid.String()
	marker = strings.Replace(marker, Separator, Void, -1)
	for i := 1; i <= 16; i++ {
		var rnd uint
		rnd, err = random.Uint(0, 255)
		if err != nil {
			return
		}
		marker = marker + hex.EncodeToString([]byte{byte(rnd)})
	}
	marker = strings.ToUpper(marker)
	return
}
