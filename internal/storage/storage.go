package storage

type Storage interface {
	GetGaugeMap() map[string]float64
	UpdateGauge(name string, value float64) error
	GetCounterMap() map[string]int64
	UpdateCounter(name string, value int64) error
}
