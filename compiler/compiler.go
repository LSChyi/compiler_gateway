package compiler

import (
	"io"
	"os"

	"compiler_gateway/compiler/ondemand"
	"compiler_gateway/compiler/stationary"

	"github.com/sirupsen/logrus"
)

const (
	FirmwareName = "ex_dactyl_v2_custom.hex"

	TypeKey = "TYPE"

	OnDemand   = "ONDEMAND"
	Stationary = "STATIONARY"
)

type Complier interface {
	Compile(keymap io.Reader) (io.Reader, error)
	Close()
}

func NewCompiler() (Complier, error) {
	compilerType, set := os.LookupEnv(TypeKey)
	if !set {
		logrus.Info("compiler type not set, use on default compiler")
		return stationary.NewCompiler()
	}

	switch compilerType {
	case Stationary:
		return stationary.NewCompiler()
	case OnDemand:
		return ondemand.NewCompiler()
	default:
		logrus.WithField("type", compilerType).Info("unimplemented, use default compiler")
		return stationary.NewCompiler()
	}
}
