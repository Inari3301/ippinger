package memstore

import (
	"github.com/Inari3301/ippinger/internel/v1/model"
	"io/ioutil"
	"os"
	"testing"
)

var (
	profiles = []model.Profile{
		{
			ID:      1,
			Name:    "aaaa",
			Surname: "bbbb",
		},
		{
			ID:      2,
			Name:    "aaaa",
			Surname: "bbbb",
		},
		{
			ID:      3,
			Name:    "aaaa",
			Surname: "bbbb",
		},
		{
			ID:      4,
			Name:    "aaaa",
			Surname: "bbbb",
		},
		{
			ID:      5,
			Name:    "aaaa",
			Surname: "bbbb",
		},
		{
			ID:      6,
			Name:    "aaaa",
			Surname: "bbbb",
		},
		{
			ID:      7,
			Name:    "aaaa",
			Surname: "bbbb",
		},
		{
			ID:      8,
			Name:    "aaaa",
			Surname: "bbbb",
		},
		{
			ID:      9,
			Name:    "aaaa",
			Surname: "bbbb",
		},
		{
			ID:      10,
			Name:    "aaaa",
			Surname: "bbbb",
		},
	}
)

func TestNewMemoryStore(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(dir)
	s, err := NewProfileStore(dir, 10, 10)
	if err != nil {
		t.Error(err)
	}

	for _, p := range profiles {
		err = s.Create(p)
		if err != nil {
			t.Error(err)
		}
	}

	s.Close()
	s = &ProfileStore{}
	s, err = NewProfileStore(dir, 10, 10)
	for _, p := range profiles {
		_, ok := s.Get(p.ID)
		if !ok {
			t.Errorf("element with id=%d not exists", p.ID)
		}
	}
}
