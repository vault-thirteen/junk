package audio

import m "github.com/vault-thirteen/SDLW/SDL/model"

type Device struct {
	Index m.Int
	Name  string
	Spec  *m.AudioSpec
}
