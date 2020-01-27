package ondemand

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
)

func genTarKeymap(source io.Reader) (io.Reader, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(source)
	data := buf.Bytes()

	buf = new(bytes.Buffer)
	tarWriter := tar.NewWriter(buf)
	header := &tar.Header{
		Name: keymapName,
		Mode: 0600,
		Size: int64(len(data)),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return nil, fmt.Errorf("error while trying to write keymap tar file header, err: %v", err)
	}
	if _, err := tarWriter.Write(data); err != nil {
		return nil, fmt.Errorf("error while trying to write keymap tar file content, err: %v", err)
	}
	if err := tarWriter.Close(); err != nil {
		return nil, fmt.Errorf("error while trying to write keymap tar file finish, err: %v", err)
	}
	return buf, nil
}

func unTarFirmware(source io.Reader) (io.Reader, error) {
	tarFile := tar.NewReader(source)
	_, err := tarFile.Next()
	if err != nil {
		return nil, err
	}
	return tarFile, nil
}
