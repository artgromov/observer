package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "", 0)

type Client struct {
	serverEndpoint string
}

func NewClient(ServerEndpoint string) *Client {
	c := new(Client)
	c.serverEndpoint = ServerEndpoint
	return c
}

func (c *Client) PushGauge(metricName string, metricValue float64) error {
	url := fmt.Sprintf("%s/update/gauge/%s/%f", c.serverEndpoint, metricName, metricValue)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		logger.Printf("failed to push %s", url)
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		logger.Printf("failed to read body %s", url) // No need to process this error somehow for now.
	}
	return nil
}

func (c *Client) PushCounter(metricName string, metricValue int64) error {
	url := fmt.Sprintf("%s/update/counter/%s/%d", c.serverEndpoint, metricName, metricValue)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		logger.Printf("failed to push %s", url)
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		logger.Printf("failed to read body %s", url) // No need to process this error somehow for now.
	}
	return nil
}
