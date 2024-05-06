package collectors

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", 0)

type Collector interface {
	Start()
	Stop()
	Poll()
	Report()
}
