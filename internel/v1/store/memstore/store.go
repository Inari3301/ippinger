package memstore

import (
	"time"
)

type Options struct {
	Path         string
	DumpInterval time.Duration
	BatchSize    uint64
}

type Store struct {
	ProfileStore
}

func New(opt Options) (*Store, error) {
	pS, err := NewProfileStore(opt.Path, opt.DumpInterval, opt.BatchSize)
	if err != nil {
		return nil, err
	}
	return &Store{
		*pS,
	}, nil
}
