package store

import "fmt"

type Updater[ValueT any] func(val ValueT) ValueT

var (
	DuplicateError = fmt.Errorf("duplicate key error")
	NotFoundError  = fmt.Errorf("key missing")
)

type Store[KeyT comparable, ValueT any] interface {
	Set(key KeyT, value ValueT) error
	Get(key KeyT) (ValueT, bool)
	Update(key KeyT, fn Updater[ValueT]) error
	Delete(key KeyT) error
}
