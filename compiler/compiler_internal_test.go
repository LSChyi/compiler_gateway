package compiler

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	testFile = "test_data/keymap.c"
)

func TestGenTarFile(t *testing.T) {
	reader := loadFile(testFile)
	res, err := genTarKeymap(reader)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

type CompilerSuite struct {
	suite.Suite
	compiler *Compiler
}

func (c *CompilerSuite) SetupTest() {
	var err error
	c.compiler, err = NewCompiler()
	if err != nil {
		panic(err)
	}
}

func (c *CompilerSuite) TearDownTest() {
	c.compiler.Close()
}

func (c *CompilerSuite) TestCreateCompiler() {
}

func (c *CompilerSuite) TestPutKeymap() {
	reader := loadFile(testFile)
	err := c.compiler.putKeymap(reader)
	c.Nil(err)
}

func (c *CompilerSuite) TestStartContainer() {
	err := c.compiler.startContainer()
	c.Nil(err)
}

func (c *CompilerSuite) TestWaitContainer() {
	reader := loadFile(testFile)
	err := c.compiler.putKeymap(reader)
	check(err)
	err = c.compiler.startContainer()
	check(err)
	err = c.compiler.waitContainer()
	c.Nil(err)
}

func (c *CompilerSuite) TestCopyFirmware() {
	reader := loadFile(testFile)
	err := c.compiler.putKeymap(reader)
	check(err)
	err = c.compiler.startContainer()
	check(err)
	err = c.compiler.waitContainer()
	check(err)
	reader, err = c.compiler.getFirmware()
	c.NotNil(reader)
	c.Nil(err)
}

func (c *CompilerSuite) TestCompile() {
	reader := loadFile(testFile)
	firmwareReader, err := c.compiler.Compile(reader)
	c.NotNil(firmwareReader)
	c.Nil(err)
}

func TestCompilerSuite(t *testing.T) {
	suite.Run(t, new(CompilerSuite))
}

func loadFile(path string) io.Reader {
	f, err := os.Open(path)
	check(err)
	return f
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
