package memstore

import (
	"fmt"
	"github.com/Inari3301/ippinger/pkg/dumper"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Inari3301/ippinger/internel/v1/model"
	"github.com/Inari3301/ippinger/pkg/store"
)

const (
	jsonDumpExtend = "ippinger.dump.json"
)

var (
	currentDumpFilename string
)

type ProfileStore struct {
	s         store.Store[uint64, model.Profile]
	c         chan uint64
	closeC    chan bool
	pathToDir string
	batchSize uint64
	timer     *time.Timer
}

func NewProfileStore(dumpPath string, dumpTimeoutSec time.Duration, batchSize uint64) (*ProfileStore, error) {
	p := &ProfileStore{
		s:         store.NewRWMapStore[uint64, model.Profile](),
		c:         make(chan uint64),
		closeC:    make(chan bool),
		pathToDir: dumpPath,
		timer:     time.NewTimer(time.Second * dumpTimeoutSec),
		batchSize: batchSize,
	}
	files, err := ioutil.ReadDir(dumpPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), jsonDumpExtend) {
			log.Printf("try get info from dump file: %s", file.Name())
			filename := path.Join(dumpPath, file.Name())
			f, err := os.Open(filename)
			if err != nil {
				log.Printf("%v", err)
				continue
			}
			d := dumper.JsonDumper{
				File: f,
			}
			err = d.FromDump(p.s)
			if err != nil {
				log.Printf("%v", err)
				_ = f.Close()
				continue
			}
			_ = f.Close()
			currentDumpFilename = filename
			break
		}
	}

	go p.dumpLoop()
	return p, nil
}

func (ps *ProfileStore) Close() {
	ps.closeC <- true
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
	if !ok {
		return p, ok
	}
	return p, ok
}

func (ps *ProfileStore) dump() error {
	dumpFileName := path.Join(ps.pathToDir, strconv.FormatInt(time.Now().Unix(), 10)+"."+jsonDumpExtend)
	t := strings.Split(dumpFileName, " ")
	dumpFileName = strings.Join(t, "_")
	f, err := os.Create(dumpFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	d := dumper.JsonDumper{
		File: f,
	}

	err = d.Dump(ps.s)
	if err != nil {
		return err
	}
	_ = os.Remove(currentDumpFilename)
	currentDumpFilename = dumpFileName
	return nil
}

func (ps *ProfileStore) dumpLoop() {
	var count uint64
	for {
		select {
		case <-ps.c:
			count++
			if count < ps.batchSize {
				continue
			}
			count = 0
			err := ps.dump()
			if err != nil {
				log.Printf("%v", err)
			}
		case <-ps.timer.C:
			if count == 0 {
				continue
			}
			err := ps.dump()
			if err != nil {
				log.Printf("%v", err)
			}
		case <-ps.closeC:
			if count == 0 {
				return
			}
			err := ps.dump()
			if err != nil {
				log.Printf("%v", err)
			}
			currentDumpFilename = ""
			return
		}
	}
}
