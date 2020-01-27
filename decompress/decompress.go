package decompress

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
)

func Unzip(source io.Reader) (*zip.Reader, error) {
	b, err := ioutil.ReadAll(source)
	if err != nil {
		return nil, err
	}
	reader, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func UnzipAndLoadTarget(source io.Reader, base string) (io.ReadCloser, error) {
	reader, err := Unzip(source)
	if err != nil {
		return nil, err
	}
	target := findByFileBase(base, reader.File)
	if target == nil {
		return nil, fmt.Errorf("can not find target file %s", base)
	}
	res, err := target.Open()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func findByFileBase(base string, files []*zip.File) *zip.File {
	for _, f := range files {
		if path.Base(f.Name) == base {
			return f
		}
	}
	return nil
}
