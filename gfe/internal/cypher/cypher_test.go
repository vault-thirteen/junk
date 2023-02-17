package cypher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//TODO: Доделать этот пакет, когда информация о Vault станет доступной.

func TestParseRsaPublicKey(t *testing.T) {
	// Arrange.
	publicKeyPEM := `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwjZtXrsyqk5zHxIG3toi
5drMsZW31rGaAVEY6E6dYsly1hNTb4kmk5J3IdhlwaJKqm0/0I0EVrVoZWwPKdgQ
xD8S3ekrMcSU4b6D6YInGOb5TrTLRkwlnBJzmfMZekngTadBb40hC+1ekQw2zln2
9e3Hmvn4hTKtY4AaG8dgiasd+ididnQQqhgZgdmJChkSvtoVcioPVGLGE9Yv6EbZ
7y/4aIWatMxrywOPoH85UHKT32XtAKtBzRLL/lvvBoeyzNCjZchQdm0fBcbC1yjI
Z3YCSdyMt+DfCKOy+BZYSHEpdX/dOyMt1rZFIroi0WpAt2xd6+z8W7rr+ru1QWrl
Tw0sBj3SM6mupsrtUOJsynY2sv6IIZ71huRmFvjEKWhrjb5A/s+XJp7dVd15noJv
UeCag3dJ+BM8Rd2TYoVy3F9sRBmhgh1v9l/eZYkTYiGBHMJu+gB3JtR3H/fVFACa
9vpkuOrm5Fy3p9xZsft0i80NtkY5Ad5e9tAA6ImhF+lp8Tkt+flJypqDnTNTGAH5
IVP283JGshzVDFOBi3xO25NiiznAXYmrfbhuxWErgnJZ+WbPCqTv7DaJ+v+d2vWx
EU7W9Yr4JH6Xtz7vRoyoZ35Nn6Fc3SWRDbvTCm675WacoEO3xkEqG4jMgxIlsWJf
Un51quvnqXEng6SDHvuXMNcCAwEAAQ==
-----END PUBLIC KEY-----
`

	// Act.
	_, err := ParseRsaPublicKey(publicKeyPEM)

	// Assert.
	assert.NoError(t, err)
}
