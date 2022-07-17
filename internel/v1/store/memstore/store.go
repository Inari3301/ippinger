package memstore

import (
	s "github.com/Inari3301/ippinger/internel/v1/store"
	"time"
)

type Options struct {
	Path         string
	DumpInterval time.Duration
	BatchSize    uint64
}

func NewMemoryStore(opt Options) *s.Store {
	return &s.Store{
		ProfileStore: NewProfileStore(opt.Path, opt.DumpInterval, opt.BatchSize),
	}
}
