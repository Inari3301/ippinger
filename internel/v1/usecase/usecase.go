package usecase

import (
	"github.com/Inari3301/ippinger/internel/v1/model"
	"time"
)

type Ping interface {
	Ping(ip string, timeout time.Duration) (model.PingResult, error)
	PingByCsv(csv []byte) (model.PingByCsvResult, error)
}

type Profile interface {
	Create(profile model.Profile) error
	Check(ID uint64) bool
}

type UseCase interface {
	Ping
	Profile
}
