package store

import "fmt"

type Updater func(val interface{}) interface{}

var (
	DuplicateError = fmt.Errorf("duplicate key error")
	NotFoundError  = fmt.Errorf("key missing")
)

type Store[T comparable] interface {
	Set(key T, value interface{}) error
	Get(key T) (interface{}, bool)
	Update(key T, fn Updater) error
	Delete(key T) error
}
