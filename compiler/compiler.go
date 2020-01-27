package compiler

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"online_compiler/compiler/ondemand"

	"github.com/sirupsen/logrus"
)

const (
	WaitCompilerTimeout = 10 * time.Second
	RetryNewCompiler    = time.Second
	FirmwareName        = "ex_dactyl_v2_custom.hex"

	TypeKey    = "TYPE"
	WorkersKey = "WORKERS"

	OnDemand = "ONDEMAND"
)

type Compiler interface {
	Compile(keymap io.Reader) (io.Reader, error)
	Close()
}

type Manager struct {
	compilerType string
	workerChan   chan Compiler
}

func NewManager() (*Manager, error) {
	res := new(Manager)

	if compilerType, set := os.LookupEnv(TypeKey); !set {
		logrus.Info("compiler type not set, use default compiler")
		res.compilerType = OnDemand
	} else {
		switch compilerType {
		case OnDemand:
			res.compilerType = OnDemand
		default:
			logrus.WithField("type", compilerType).Info("unimplemented, use default compiler")
			res.compilerType = OnDemand
		}
	}

	workers := runtime.NumCPU()
	res.workerChan = make(chan Compiler, workers)
	for i := 0; i < workers; i++ {
		if compiler, err := res.newCompiler(); err != nil {
			defer res.Close()
			return nil, err
		} else {
			res.workerChan <- compiler
		}
	}
	return res, nil
}

func (m *Manager) GetCompiler() (Compiler, error) {
	ctx, cancel := context.WithTimeout(context.Background(), WaitCompilerTimeout)
	defer cancel()
	select {
	case compiler := <-m.workerChan:
		return compiler, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout when getting compiler")
	}
	return nil, nil
}

func (m *Manager) newCompiler() (Compiler, error) {
	switch m.compilerType {
	case OnDemand:
		return ondemand.NewCompiler()
	default:
		return ondemand.NewCompiler()
	}
}

func (m *Manager) ReturnCompiler(c Compiler) {
	c.Close()
	var compiler Compiler
	var err error
	for compiler, err = m.newCompiler(); err != nil; {
		logrus.WithError(err).Error("can not get new compiler, retry in", RetryNewCompiler)
		<-time.After(RetryNewCompiler)
	}
	m.workerChan <- compiler
}

func (m *Manager) Close() {
	close(m.workerChan)
	for compiler := range m.workerChan {
		compiler.Close()
	}
}
