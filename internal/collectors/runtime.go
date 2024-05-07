package collectors

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type RuntimeCollector struct {
	ServerEndpoint string

	gaugeMap   map[string]float64
	counterMap map[string]int64
	lock       sync.Mutex

	pollTicker   time.Ticker
	pollStop     chan bool
	reportTicker time.Ticker
	reportStop   chan bool
}

func NewRuntimeCollector(ServerEndpoint string, PollInterval time.Duration, ReportInterval time.Duration) *RuntimeCollector {
	rc := new(RuntimeCollector)
	rc.ServerEndpoint = ServerEndpoint
	rc.gaugeMap = make(map[string]float64)
	rc.counterMap = make(map[string]int64)
	rc.pollTicker = *time.NewTicker(PollInterval)
	rc.pollStop = make(chan bool, 1)
	rc.reportTicker = *time.NewTicker(ReportInterval)
	rc.reportStop = make(chan bool, 1)
	return rc
}

func (cl *RuntimeCollector) Start() {
	logger.Printf("RuntimeCollector starting")
	go cl.Poll()
	go cl.Report()
	logger.Printf("RuntimeCollector started")
}

func (cl *RuntimeCollector) Stop() {
	logger.Printf("RuntimeCollector stopping")
	cl.pollTicker.Stop()
	cl.reportTicker.Stop()
	cl.pollStop <- true
	cl.reportStop <- true
	<-cl.pollStop
	<-cl.reportStop
	logger.Printf("RuntimeCollector stopped")
}

func (cl *RuntimeCollector) Poll() {
	logger.Printf("RuntimeCollector poll goroutine started")
	defer func() { cl.pollStop <- true }()
	for {
		select {
		case <-cl.pollTicker.C:
			func() {
				cl.lock.Lock()
				defer cl.lock.Unlock()
				logger.Printf("RuntimeCollector poll iteration started")

				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				cl.gaugeMap["Alloc"] = float64(memStats.Alloc)
				cl.gaugeMap["BuckHashSys"] = float64(memStats.BuckHashSys)
				cl.gaugeMap["Frees"] = float64(memStats.Frees)
				cl.gaugeMap["GCCPUFraction"] = float64(memStats.GCCPUFraction)
				cl.gaugeMap["GCSys"] = float64(memStats.GCSys)
				cl.gaugeMap["HeapAlloc"] = float64(memStats.HeapAlloc)
				cl.gaugeMap["HeapIdle"] = float64(memStats.HeapIdle)
				cl.gaugeMap["HeapInuse"] = float64(memStats.HeapInuse)
				cl.gaugeMap["HeapObjects"] = float64(memStats.HeapObjects)
				cl.gaugeMap["HeapReleased"] = float64(memStats.HeapReleased)
				cl.gaugeMap["HeapSys"] = float64(memStats.HeapSys)
				cl.gaugeMap["LastGC"] = float64(memStats.LastGC)
				cl.gaugeMap["Lookups"] = float64(memStats.Lookups)
				cl.gaugeMap["MCacheInuse"] = float64(memStats.MCacheInuse)
				cl.gaugeMap["MCacheSys"] = float64(memStats.MCacheSys)
				cl.gaugeMap["MSpanInuse"] = float64(memStats.MSpanInuse)
				cl.gaugeMap["MSpanSys"] = float64(memStats.MSpanSys)
				cl.gaugeMap["Mallocs"] = float64(memStats.Mallocs)
				cl.gaugeMap["NextGC"] = float64(memStats.NextGC)
				cl.gaugeMap["NumForcedGC"] = float64(memStats.NumForcedGC)
				cl.gaugeMap["NumGC"] = float64(memStats.NumGC)
				cl.gaugeMap["OtherSys"] = float64(memStats.OtherSys)
				cl.gaugeMap["PauseTotalNs"] = float64(memStats.PauseTotalNs)
				cl.gaugeMap["StackInuse"] = float64(memStats.StackInuse)
				cl.gaugeMap["StackSys"] = float64(memStats.StackSys)
				cl.gaugeMap["Sys"] = float64(memStats.Sys)
				cl.gaugeMap["TotalAlloc"] = float64(memStats.TotalAlloc)
				cl.gaugeMap["RandomValue"] = rand.Float64()

				cl.counterMap["PollCount"] += 1

				logger.Printf("RuntimeCollector poll iteration finished")
			}()

		case <-cl.pollStop:
			logger.Printf("RuntimeCollector poll goroutine exiting by signal")
			return
		}
	}
}

func (cl *RuntimeCollector) Report() {
	logger.Printf("RuntimeCollector report goroutine started")
	defer func() { cl.reportStop <- true }()
	for {
		select {
		case <-cl.reportTicker.C:
			func() {
				cl.lock.Lock()
				defer cl.lock.Unlock()
				logger.Printf("RuntimeCollector report iteration started")
				for metricName, metricValue := range cl.gaugeMap {
					url := fmt.Sprintf("%s/update/gauge/%s/%f", cl.ServerEndpoint, metricName, metricValue)
					resp, err := http.Post(url, "text/plain", nil)
					if err != nil {
						logger.Printf("failed to push %s", url)
					}
					defer resp.Body.Close()
					_, err = io.Copy(io.Discard, resp.Body)
					if err != nil {
						logger.Printf("failed to read body %s", url)
					}
				}
				for metricName, metricValue := range cl.counterMap {
					url := fmt.Sprintf("%s/update/counter/%s/%d", cl.ServerEndpoint, metricName, metricValue)
					resp, err := http.Post(url, "text/plain", nil)
					if err != nil {
						logger.Printf("failed to push %s", url)
					}
					defer resp.Body.Close()
					_, err = io.Copy(io.Discard, resp.Body)
					if err != nil {
						logger.Printf("failed to read body %s", url)
					}
				}

				logger.Printf("RuntimeCollector report iteration finished")
			}()

		case <-cl.reportStop:
			logger.Printf("RuntimeCollector report goroutine exiting by signal")
			return
		}
	}
}
