package memstore

import (
	"fmt"
	"github.com/Inari3301/ippinger/pkg/dumper"
	"log"
	"os"
	"path"
	"time"

	"github.com/Inari3301/ippinger/internel/v1/model"
	"github.com/Inari3301/ippinger/pkg/store"
)

const (
	dateTimeLayout = "2006-02-01"
)

type ProfileStore struct {
	s         store.Store[uint64]
	c         chan uint64
	pathToDir string
	batchSize uint64
	timer     *time.Timer
}

func NewProfileStore(path string, dumpTimeoutSec time.Duration, batchSize uint64) *ProfileStore {
	p := &ProfileStore{
		s:         store.NewRWMapStore[uint64](),
		c:         make(chan uint64),
		pathToDir: path,
		timer:     time.NewTimer(time.Second * dumpTimeoutSec),
		batchSize: batchSize,
	}

	go p.dumpLoop()
	return p
}

func (ps *ProfileStore) Create(p model.Profile) error {
	err := ps.s.Set(p.ID, p)
	if err != nil {
		return fmt.Errorf("error while profile create %w", err)
	}

	ps.c <- p.ID

	return nil
}

func (ps *ProfileStore) Get(ID uint64) (model.Profile, bool) {
	p, ok := ps.s.Get(ID)
	return p.(model.Profile), ok
}

func (ps *ProfileStore) dump() (string, error) {
	dumpFileName := path.Join(ps.pathToDir, time.Now().Format(dateTimeLayout)+".ippinger.dump.json")
	f, err := os.Open(dumpFileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	d := dumper.JsonDumper{
		File: f,
	}

	err = d.Dump(ps.s)
	return dumpFileName, err
}

func (ps *ProfileStore) dumpLoop() {
	batch := make([]uint64, ps.batchSize)
	i := uint64(0)
	var dumpFileName string
	for {
		select {
		case ID := <-ps.c:
			if i >= ps.batchSize {
				var err error
				newDumpFileName, err := ps.dump()
				if err != nil {
					i = 0
					log.Printf("%v", err)
					continue
				}
				_ = os.Remove(dumpFileName)
				dumpFileName = newDumpFileName
				for _, updatedID := range batch {
					log.Printf("updates associated with user id=%d are flushed to disk", updatedID)
				}
				i = 0
			}
			batch[i] = ID
			i++
		case <-ps.timer.C:
			if i == 0 {
				continue
			}

			newDumpFileName, err := ps.dump()
			i = 0
			if err != nil {
				i = 0
				log.Printf("%v", err)
				continue
			}
			_ = os.Remove(dumpFileName)
			dumpFileName = newDumpFileName
		}
	}
}
