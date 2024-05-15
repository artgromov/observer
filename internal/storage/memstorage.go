package storage

import (
	"fmt"
	"sort"
	"sync"
)

type MemStorage struct {
	gaugeMap   map[string]float64
	counterMap map[string]int64
	lock       sync.Mutex
}

func NewMemStorage() *MemStorage {
	s := new(MemStorage)
	s.gaugeMap = make(map[string]float64)
	s.counterMap = make(map[string]int64)
	return s
}

func (s *MemStorage) GetGauge(name string) (float64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	value, ok := s.gaugeMap[name]
	if ok {
		return value, nil
	}
	return 0, fmt.Errorf("gauge \"%s\" not found", name)
}

func (s *MemStorage) UpdateGauge(name string, value float64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.gaugeMap[name] = value
	return nil
}

func (s *MemStorage) GetCounter(name string) (int64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	value, ok := s.counterMap[name]
	if ok {
		return value, nil
	}
	return 0, fmt.Errorf("counter \"%s\" not found", name)
}

func (s *MemStorage) UpdateCounter(name string, value int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.counterMap[name]
	if ok {
		s.counterMap[name] += value
	} else {
		s.counterMap[name] = value
	}
	return nil
}

func (s *MemStorage) Dump() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	var result []string

	gaugeKeys := make([]string, 0)
	for k := range s.gaugeMap {
		gaugeKeys = append(gaugeKeys, k)
	}
	sort.Strings(gaugeKeys)

	for _, k := range gaugeKeys {
		v := s.gaugeMap[k]
		result = append(result, fmt.Sprintf("gauge %s %f", k, v))
	}

	counterKeys := make([]string, 0)
	for k := range s.counterMap {
		counterKeys = append(counterKeys, k)
	}
	sort.Strings(counterKeys)

	for _, k := range counterKeys {
		v := s.counterMap[k]
		result = append(result, fmt.Sprintf("counter %s %d", k, v))
	}
	result = append(result, "")
	return result
}
