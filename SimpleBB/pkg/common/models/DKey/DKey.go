package dk

import (
	"bytes"
	"encoding/hex"
	"errors"
	"strings"
	"sync/atomic"

	"github.com/vault-thirteen/auxie/random"
)

const (
	ErrKeySize = "key size is wrong"
)

type DKey struct {
	size       int
	bytes      []byte
	str        string
	readsCount atomic.Int32
}

func NewDKey(size int) (key *DKey, err error) {
	if size <= 0 {
		return nil, errors.New(ErrKeySize)
	}

	key = &DKey{
		size: size,
	}

	key.bytes, err = random.GenerateRandomBytes(key.size)
	if err != nil {
		return nil, err
	}

	key.str = strings.ToUpper(hex.EncodeToString(key.bytes))

	key.readsCount = atomic.Int32{}
	key.readsCount.Store(0)

	return key, nil
}

func (k *DKey) GetBytes() []byte {
	ok := k.readsCount.CompareAndSwap(0, 1)
	if !ok {
		return nil
	}

	return k.bytes
}

func (k *DKey) GetString() string {
	ok := k.readsCount.CompareAndSwap(0, 1)
	if !ok {
		return ""
	}

	return k.str
}

func (k *DKey) CheckBytes(x []byte) bool {
	return bytes.Compare(x, k.bytes) == 0
}

func (k *DKey) CheckString(x string) bool {
	return x == k.str
}
