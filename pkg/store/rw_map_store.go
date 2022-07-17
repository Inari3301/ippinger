package store

import (
	"encoding/json"
	"sync"
)

type RWMapStore[KeyT comparable, ValueT any] struct {
	data map[KeyT]ValueT
	lock sync.RWMutex
}

func NewRWMapStore[KeyT comparable, ValueT any]() *RWMapStore[KeyT, ValueT] {
	return &RWMapStore[KeyT, ValueT]{
		data: make(map[KeyT]ValueT),
	}
}

func (s *RWMapStore[KeyT, ValueT]) Set(key KeyT, value ValueT) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.data[key]
	if ok {
		return DuplicateError
	}
	s.data[key] = value
	return nil
}

func (s *RWMapStore[KeyT, ValueT]) Get(key KeyT) (ValueT, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	d, ok := s.data[key]
	return d, ok
}

func (s *RWMapStore[KeyT, ValueT]) Update(key KeyT, fn Updater[ValueT]) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	d, ok := s.data[key]
	if !ok {
		return NotFoundError
	}

	s.data[key] = fn(d)
	return nil
}

func (s *RWMapStore[KeyT, ValueT]) Delete(key KeyT) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.data[key]
	if !ok {
		return NotFoundError
	}
	delete(s.data, key)
	return nil
}

func (s *RWMapStore[KeyT, ValueT]) MarshalJSON() ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	b, err := json.Marshal(s.data)
	return b, err
}

func (s *RWMapStore[KeyT, ValueT]) UnmarshalJSON(data []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return json.Unmarshal(data, &s.data)
}
