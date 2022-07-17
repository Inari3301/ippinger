package store

import (
	"encoding/json"
	"sync"
)

type RWMapStore[T comparable] struct {
	data map[T]interface{}
	lock sync.RWMutex
}

func NewRWMapStore[T comparable]() *RWMapStore[T] {
	return &RWMapStore[T]{
		data: make(map[T]interface{}),
	}
}

func (s *RWMapStore[T]) Set(key T, value interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.data[key]
	if ok {
		return DuplicateError
	}
	s.data[key] = value
	return nil
}

func (s *RWMapStore[T]) Get(key T) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	d, ok := s.data[key]
	return d, ok
}

func (s *RWMapStore[T]) Update(key T, fn Updater) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	d, ok := s.data[key]
	if !ok {
		return NotFoundError
	}

	s.data[key] = fn(d)
	return nil
}

func (s *RWMapStore[T]) Delete(key T) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.data[key]
	if !ok {
		return NotFoundError
	}
	delete(s.data, key)
	return nil
}

func (s *RWMapStore[T]) MarshalJSON() ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	b, err := json.Marshal(s.data)
	return b, err
}

func (s *RWMapStore[T]) UnmarshalJSON(data []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return json.Unmarshal(data, &s.data)
}
