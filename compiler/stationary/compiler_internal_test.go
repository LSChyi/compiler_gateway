package stationary

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	keymapPath = "test_data/keymap.c"
	testUrl    = "http://localhost:8080/"
)

func TestSendRequest(t *testing.T) {
	f := loadKeymap(keymapPath)
	defer f.Close()
	c, err := NewCompiler()
	check(err)
	defer c.Close()
	firmware, err := c.sendRequest(testUrl, f)
	assert.NotNil(t, firmware)
	assert.Nil(t, err)
}

func loadKeymap(path string) *os.File {
	f, err := os.Open(path)
	check(err)
	return f
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
