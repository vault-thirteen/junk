package jwt

import (
	"crypto/rsa"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewJwt(t *testing.T) {
	var (
		err error
		jwt *Jwt
	)

	// Arrange.
	err = os.Setenv("GFE_JWT_KEY_SOURCE_TYPE", "1")
	assert.NoError(t, err)

	err = os.Setenv("GFE_JWT_KEY_DSN", "")
	assert.NoError(t, err)

	err = os.Setenv("GFE_JWT_KEY_VALUE", `-----BEGIN PUBLIC KEY-----.MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwjZtXrsyqk5zHxIG3toi.5drMsZW31rGaAVEY6E6dYsly1hNTb4kmk5J3IdhlwaJKqm0/0I0EVrVoZWwPKdgQ.xD8S3ekrMcSU4b6D6YInGOb5TrTLRkwlnBJzmfMZekngTadBb40hC+1ekQw2zln2.9e3Hmvn4hTKtY4AaG8dgiasd+ididnQQqhgZgdmJChkSvtoVcioPVGLGE9Yv6EbZ.7y/4aIWatMxrywOPoH85UHKT32XtAKtBzRLL/lvvBoeyzNCjZchQdm0fBcbC1yjI.Z3YCSdyMt+DfCKOy+BZYSHEpdX/dOyMt1rZFIroi0WpAt2xd6+z8W7rr+ru1QWrl.Tw0sBj3SM6mupsrtUOJsynY2sv6IIZ71huRmFvjEKWhrjb5A/s+XJp7dVd15noJv.UeCag3dJ+BM8Rd2TYoVy3F9sRBmhgh1v9l/eZYkTYiGBHMJu+gB3JtR3H/fVFACa.9vpkuOrm5Fy3p9xZsft0i80NtkY5Ad5e9tAA6ImhF+lp8Tkt+flJypqDnTNTGAH5.IVP283JGshzVDFOBi3xO25NiiznAXYmrfbhuxWErgnJZ+WbPCqTv7DaJ+v+d2vWx.EU7W9Yr4JH6Xtz7vRoyoZ35Nn6Fc3SWRDbvTCm675WacoEO3xkEqG4jMgxIlsWJf.Un51quvnqXEng6SDHvuXMNcCAwEAAQ==.-----END PUBLIC KEY-----.`)
	assert.NoError(t, err)

	var logger zerolog.Logger
	logger = zerolog.New(os.Stderr).With().Timestamp().Stack().Logger()

	// Act.
	jwt, err = NewJwt(&logger)

	// Assert.
	assert.NoError(t, err)
	assert.NotEqual(t, nil, jwt)

	// Очищаем О.С. от мусора после теста.
	err = os.Setenv("GFE_JWT_KEY_SOURCE_TYPE", "")
	assert.NoError(t, err)
	err = os.Setenv("GFE_JWT_KEY_DSN", "")
	assert.NoError(t, err)
	err = os.Setenv("GFE_JWT_KEY_VALUE", "")
	assert.NoError(t, err)
}

func TestJwt_GetRsaPublicKey(t *testing.T) {
	var (
		err error
		jwt *Jwt
		key *rsa.PublicKey
	)

	// Arrange.
	err = os.Setenv("GFE_JWT_KEY_SOURCE_TYPE", "1")
	assert.NoError(t, err)

	err = os.Setenv("GFE_JWT_KEY_DSN", "")
	assert.NoError(t, err)

	err = os.Setenv("GFE_JWT_KEY_VALUE", `-----BEGIN PUBLIC KEY-----.MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwjZtXrsyqk5zHxIG3toi.5drMsZW31rGaAVEY6E6dYsly1hNTb4kmk5J3IdhlwaJKqm0/0I0EVrVoZWwPKdgQ.xD8S3ekrMcSU4b6D6YInGOb5TrTLRkwlnBJzmfMZekngTadBb40hC+1ekQw2zln2.9e3Hmvn4hTKtY4AaG8dgiasd+ididnQQqhgZgdmJChkSvtoVcioPVGLGE9Yv6EbZ.7y/4aIWatMxrywOPoH85UHKT32XtAKtBzRLL/lvvBoeyzNCjZchQdm0fBcbC1yjI.Z3YCSdyMt+DfCKOy+BZYSHEpdX/dOyMt1rZFIroi0WpAt2xd6+z8W7rr+ru1QWrl.Tw0sBj3SM6mupsrtUOJsynY2sv6IIZ71huRmFvjEKWhrjb5A/s+XJp7dVd15noJv.UeCag3dJ+BM8Rd2TYoVy3F9sRBmhgh1v9l/eZYkTYiGBHMJu+gB3JtR3H/fVFACa.9vpkuOrm5Fy3p9xZsft0i80NtkY5Ad5e9tAA6ImhF+lp8Tkt+flJypqDnTNTGAH5.IVP283JGshzVDFOBi3xO25NiiznAXYmrfbhuxWErgnJZ+WbPCqTv7DaJ+v+d2vWx.EU7W9Yr4JH6Xtz7vRoyoZ35Nn6Fc3SWRDbvTCm675WacoEO3xkEqG4jMgxIlsWJf.Un51quvnqXEng6SDHvuXMNcCAwEAAQ==.-----END PUBLIC KEY-----.`)
	assert.NoError(t, err)

	var logger zerolog.Logger
	logger = zerolog.New(os.Stderr).With().Timestamp().Stack().Logger()

	jwt, err = NewJwt(&logger)
	assert.NoError(t, err)

	// Act.
	key = jwt.GetRsaPublicKey()

	// Assert.
	assert.NotEqual(t, nil, key)

	// Очищаем О.С. от мусора после теста.
	err = os.Setenv("GFE_JWT_KEY_SOURCE_TYPE", "")
	assert.NoError(t, err)
	err = os.Setenv("GFE_JWT_KEY_DSN", "")
	assert.NoError(t, err)
	err = os.Setenv("GFE_JWT_KEY_VALUE", "")
	assert.NoError(t, err)
}
