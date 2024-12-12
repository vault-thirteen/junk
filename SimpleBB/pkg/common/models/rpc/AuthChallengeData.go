package rpc

type AuthChallengeData []byte

func (acd AuthChallengeData) Bytes() []byte {
	return acd
}
