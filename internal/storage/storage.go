package storage

type Storage interface {
	GetGauge(name string) (float64, error)
	UpdateGauge(name string, value float64) error

	GetCounter(name string) (int64, error)
	UpdateCounter(name string, value int64) error

	Dump() []string
}
