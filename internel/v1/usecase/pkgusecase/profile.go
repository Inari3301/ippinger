package pkgusecase

import (
	"fmt"
	"github.com/Inari3301/ippinger/internel/v1/model"
	"github.com/Inari3301/ippinger/internel/v1/store"
)

type Profile struct {
	s store.Store
}

func (p *Profile) Create(profile model.Profile) error {
	err := p.s.Create(profile)
	if err != nil {
		return fmt.Errorf("error while create profile with id=%d: %v", profile.ID, err)
	}

	return nil
}

func (p *Profile) Check(ID uint64) bool {
	_, ok := p.s.Get(ID)
	return ok
}
