package collectors

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/artgromov/observer/internal/client"
	"github.com/artgromov/observer/internal/logger"
	"go.uber.org/zap"
)

type RuntimeCollector struct {
	client *client.Client
	l      *zap.Logger

	gaugeMap   map[string]float64
	counterMap map[string]int64
	lock       sync.Mutex

	pollTicker   time.Ticker
	pollStop     chan bool
	reportTicker time.Ticker
	reportStop   chan bool
}

func NewRuntimeCollector(Client *client.Client, PollInterval time.Duration, ReportInterval time.Duration) *RuntimeCollector {
	rc := new(RuntimeCollector)
	rc.client = Client
	rc.l = logger.Get()
	rc.gaugeMap = make(map[string]float64)
	rc.counterMap = make(map[string]int64)
	rc.pollTicker = *time.NewTicker(PollInterval)
	rc.pollStop = make(chan bool, 1)
	rc.reportTicker = *time.NewTicker(ReportInterval)
	rc.reportStop = make(chan bool, 1)
	return rc
}

func (rc *RuntimeCollector) Start() {
	rc.l.Info("RuntimeCollector starting")
	go rc.Poll()
	go rc.Report()
	rc.l.Info("RuntimeCollector started")
}

func (rc *RuntimeCollector) Stop() {
	rc.l.Info("RuntimeCollector stopping")
	rc.pollTicker.Stop()
	rc.reportTicker.Stop()
	rc.pollStop <- true
	rc.reportStop <- true
	<-rc.pollStop
	<-rc.reportStop
	rc.l.Info("RuntimeCollector stopped")
}

func (rc *RuntimeCollector) Poll() {
	rc.l.Info("RuntimeCollector poll goroutine started")
	defer func() { rc.pollStop <- true }()
	for {
		select {
		case <-rc.pollTicker.C:
			func() {
				rc.lock.Lock()
				defer rc.lock.Unlock()
				rc.l.Info("RuntimeCollector poll iteration started")

				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				rc.gaugeMap["Alloc"] = float64(memStats.Alloc)
				rc.gaugeMap["BuckHashSys"] = float64(memStats.BuckHashSys)
				rc.gaugeMap["Frees"] = float64(memStats.Frees)
				rc.gaugeMap["GCCPUFraction"] = float64(memStats.GCCPUFraction)
				rc.gaugeMap["GCSys"] = float64(memStats.GCSys)
				rc.gaugeMap["HeapAlloc"] = float64(memStats.HeapAlloc)
				rc.gaugeMap["HeapIdle"] = float64(memStats.HeapIdle)
				rc.gaugeMap["HeapInuse"] = float64(memStats.HeapInuse)
				rc.gaugeMap["HeapObjects"] = float64(memStats.HeapObjects)
				rc.gaugeMap["HeapReleased"] = float64(memStats.HeapReleased)
				rc.gaugeMap["HeapSys"] = float64(memStats.HeapSys)
				rc.gaugeMap["LastGC"] = float64(memStats.LastGC)
				rc.gaugeMap["Lookups"] = float64(memStats.Lookups)
				rc.gaugeMap["MCacheInuse"] = float64(memStats.MCacheInuse)
				rc.gaugeMap["MCacheSys"] = float64(memStats.MCacheSys)
				rc.gaugeMap["MSpanInuse"] = float64(memStats.MSpanInuse)
				rc.gaugeMap["MSpanSys"] = float64(memStats.MSpanSys)
				rc.gaugeMap["Mallocs"] = float64(memStats.Mallocs)
				rc.gaugeMap["NextGC"] = float64(memStats.NextGC)
				rc.gaugeMap["NumForcedGC"] = float64(memStats.NumForcedGC)
				rc.gaugeMap["NumGC"] = float64(memStats.NumGC)
				rc.gaugeMap["OtherSys"] = float64(memStats.OtherSys)
				rc.gaugeMap["PauseTotalNs"] = float64(memStats.PauseTotalNs)
				rc.gaugeMap["StackInuse"] = float64(memStats.StackInuse)
				rc.gaugeMap["StackSys"] = float64(memStats.StackSys)
				rc.gaugeMap["Sys"] = float64(memStats.Sys)
				rc.gaugeMap["TotalAlloc"] = float64(memStats.TotalAlloc)
				rc.gaugeMap["RandomValue"] = rand.Float64()

				rc.counterMap["PollCount"] += 1

				rc.l.Info("RuntimeCollector poll iteration finished")
			}()

		case <-rc.pollStop:
			rc.l.Info("RuntimeCollector poll goroutine exiting by signal")
			return
		}
	}
}

func (rc *RuntimeCollector) Report() {
	rc.l.Info("RuntimeCollector report goroutine started")
	defer func() { rc.reportStop <- true }()
	for {
		select {
		case <-rc.reportTicker.C:
			func() {
				rc.lock.Lock()
				defer rc.lock.Unlock()
				rc.l.Info("RuntimeCollector report iteration started")
				for metricName, metricValue := range rc.gaugeMap {
					err := rc.client.PushGauge(metricName, metricValue)
					if err != nil {
						continue
					}
				}
				for metricName, metricValue := range rc.counterMap {
					err := rc.client.PushCounter(metricName, metricValue)
					if err != nil {
						continue
					}
					rc.counterMap[metricName] = 0 // resetting counter after successful push
				}

				rc.l.Info("RuntimeCollector report iteration finished")
			}()

		case <-rc.reportStop:
			rc.l.Info("RuntimeCollector report goroutine exiting by signal")
			return
		}
	}
}
