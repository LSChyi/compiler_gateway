package compiler

import (
	"io"
	"os"

	"online_compiler/compiler/ondemand"

	"github.com/sirupsen/logrus"
)

const (
	FirmwareName = "ex_dactyl_v2_custom.hex"

	TypeKey = "TYPE"

	OnDemand = "ONDEMAND"
)

type Compiler interface {
	Compile(keymap io.Reader) (io.Reader, error)
	Close()
}

type Factory struct {
	compilerType string
}

func NewFactory() *Factory {
	res := new(Factory)
	compilerType, set := os.LookupEnv(TypeKey)
	if !set {
		logrus.Info("compiler type not set, use on default compiler")
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
	return res
}

func (f *Factory) NewCompiler() (Compiler, error) {
	switch f.compilerType {
	case OnDemand:
		return ondemand.NewCompiler()
	default:
		return ondemand.NewCompiler()
	}
}
