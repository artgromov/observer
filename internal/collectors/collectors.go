package collectors

type Collector interface {
	Start()
	Stop()
	Poll()
	Report()
}
