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

type Complier interface {
	Compile(keymap io.Reader) (io.Reader, error)
	Close()
}

func NewCompiler() (Complier, error) {
	compilerType, set := os.LookupEnv(TypeKey)
	if !set {
		logrus.Info("compiler type not set, use on default compiler")
		return ondemand.NewCompiler()
	}

	switch compilerType {
	case OnDemand:
		return ondemand.NewCompiler()
	default:
		logrus.WithField("type", compilerType).Info("unimplemented, use default compiler")
		return ondemand.NewCompiler()
	}
}
