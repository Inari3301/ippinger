package store

import "github.com/Inari3301/ippinger/internel/v1/model"

type ProfileStore interface {
	Create(p model.Profile) error
	Get(ID uint64) (model.Profile, bool)
}

type Store struct {
	ProfileStore ProfileStore
}
