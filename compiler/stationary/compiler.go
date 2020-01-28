package stationary

import (
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	fileField       = "keymap"
	fileName        = "keymap.c"
	defaultEndpoint = "http://localhost:8080/"
	EndpointKey     = "ENDPOINT"
)

type Compiler struct {
	client             *http.Client
	firmwareReadCloser io.ReadCloser
	endpoint           string
}

func NewCompiler() (*Compiler, error) {
	endpoint := defaultEndpoint
	if endpoint, set := os.LookupEnv(EndpointKey); set {
		endpoint = endpoint
	} else {
		logrus.Info("endpoint not set, use default endpoint")
	}
	return &Compiler{
		client:   new(http.Client),
		endpoint: endpoint,
	}, nil
}

func (c *Compiler) Compile(keymap io.Reader) (io.Reader, error) {
	return c.sendRequest(c.endpoint, keymap)
}

func (c *Compiler) Close() {
	if c.firmwareReadCloser != nil {
		c.firmwareReadCloser.Close()
	}
}

func (c *Compiler) sendRequest(url string, keymap io.Reader) (io.Reader, error) {
	res, err := c.client.Post(url, "application/octet-stream", keymap)
	if err != nil {
		return nil, err
	}
	c.firmwareReadCloser = res.Body
	return c.firmwareReadCloser, nil
}
