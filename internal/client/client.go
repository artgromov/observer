package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/artgromov/observer/internal/logger"
	"go.uber.org/zap"
)

type Client struct {
	serverEndpoint string
	l              *zap.Logger
}

func NewClient(ServerEndpoint string) *Client {
	c := new(Client)
	c.serverEndpoint = ServerEndpoint
	c.l = logger.Get()
	return c
}

func (c *Client) PushGauge(metricName string, metricValue float64) error {
	url := fmt.Sprintf("%s/update/gauge/%s/%f", c.serverEndpoint, metricName, metricValue)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		c.l.Error("failed to push", zap.String("url", url))
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		c.l.Error("failed to read body", zap.String("url", url)) // No need to process this error somehow for now.
	}
	return nil
}

func (c *Client) PushCounter(metricName string, metricValue int64) error {
	url := fmt.Sprintf("%s/update/counter/%s/%d", c.serverEndpoint, metricName, metricValue)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		c.l.Error("failed to push", zap.String("url", url))
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		c.l.Error("failed to read body", zap.String("url", url)) // No need to process this error somehow for now.
	}
	return nil
}
