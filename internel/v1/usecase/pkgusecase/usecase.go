package pkgusecase

import "github.com/Inari3301/ippinger/internel/v1/store"

type UseCase struct {
	Ping
	Profile
}

func New(s store.Store) *UseCase {
	return &UseCase{
		Ping:    Ping{},
		Profile: Profile{s},
	}
}
