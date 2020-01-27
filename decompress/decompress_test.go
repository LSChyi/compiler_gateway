package decompress_test

import (
	"io"
	"os"
	"testing"

	"online_compiler/decompress"

	"github.com/stretchr/testify/assert"
)

const (
	testFile = "ergodox_ez.zip"
)

var (
	decompressedFiles = []string{
		"ergodox_ez_ergodox-ez-default-layout_source/",
		"ergodox_ez_ergodox-ez-default-layout_source/config.h",
		"ergodox_ez_ergodox-ez-default-layout_source/keymap.c",
		"ergodox_ez_ergodox-ez-default-layout_source/rules.mk",
		"ergodox_ez_ergodox-ez-default-layout_Qaazm_DZ6gK.md5",
		"README.md",
		"build.log",
		"ergodox_ez_ergodox-ez-default-layout_Qaazm_DZ6gK.hex",
	}
	nonExistFile  = "nonExist.txt"
	existFileBase = "keymap.c"
)

func TestUnzip(t *testing.T) {
	expectedFileSet := generateExpectedFileSet(decompressedFiles)

	source := loadFile(testFile)
	reader, err := decompress.Unzip(source)
	assert.Nil(t, err)
	for _, f := range reader.File {
		assert.True(t, expectedFileSet[f.Name])
		delete(expectedFileSet, f.Name)
	}
	assert.Empty(t, len(expectedFileSet))
}

func TestUnzipAndLoadWithNonExist(t *testing.T) {
	source := loadFile(testFile)
	res, err := decompress.UnzipAndLoadTarget(source, nonExistFile)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestUnzipAndLoadWithFile(t *testing.T) {
	source := loadFile(testFile)
	res, err := decompress.UnzipAndLoadTarget(source, existFileBase)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func loadFile(path string) io.Reader {
	f, err := os.Open(path)
	check(err)
	return f
}

func generateExpectedFileSet(names []string) map[string]bool {
	res := make(map[string]bool)
	for _, name := range names {
		res[name] = true
	}
	return res
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
