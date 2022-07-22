package pkgusecase

import "github.com/Inari3301/ippinger/internel/v1/store"

type UseCase struct {
	P
	Profile
}

func New(s store.Store) *UseCase {
	return &UseCase{
		P:       P{},
		Profile: Profile{s},
	}
}
