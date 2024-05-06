package storage

type MemStorage struct {
	gaugeMap   map[string]float64
	counterMap map[string]int64
}

func NewMemStorage() *MemStorage {
	s := new(MemStorage)
	s.gaugeMap = make(map[string]float64)
	s.counterMap = make(map[string]int64)
	return s
}

func (s *MemStorage) GetGaugeMap() map[string]float64 {
	return s.gaugeMap
}

func (s *MemStorage) UpdateGauge(name string, value float64) error {
	s.gaugeMap[name] = value
	return nil
}

func (s *MemStorage) GetCounterMap() map[string]int64 {
	return s.counterMap
}

func (s *MemStorage) UpdateCounter(name string, value int64) error {
	_, ok := s.counterMap[name]
	if ok {
		s.counterMap[name] += value
	} else {
		s.counterMap[name] = value
	}
	return nil
}
