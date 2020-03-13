package main

import (
	"io"
	"net/http"

	"compiler_gateway/compiler"
	"compiler_gateway/decompress"

	"github.com/sirupsen/logrus"
)

func main() {
	s := NewServer(":8080")
	defer s.Close()

	logrus.WithField("address", s.Addr).Info("running http server")
	if err := s.ListenAndServe(); err != nil {
		logrus.WithError(err).Error("run http server with error")
	}
}

func NewServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRequest)
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
func handleRequest(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("request", r).Info("handling request")
	if r.Method != http.MethodPost {
		logrus.Info("only handle post method, drop it")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	r.ParseMultipartForm(64 * 1024)
	keymapZip, _, err := r.FormFile("keymapZip")
	if err != nil {
		logrus.WithError(err).Error("can not get keymap zip")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer keymapZip.Close()

	keymap, err := decompress.UnzipAndLoadTarget(keymapZip, "keymap.c")
	if err != nil {
		logrus.WithError(err).Error("can not get keymap file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = writeCompileResult(w, keymap)
	if err != nil {
		logrus.WithError(err).Error("get error while compiling")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
}

func writeCompileResult(w http.ResponseWriter, keymap io.Reader) error {
	c, err := compiler.NewCompiler()
	if err != nil {
		logrus.WithError(err).Error("can not create compiler")
		return err
	}
	defer c.Close()

	firmware, err := c.Compile(keymap)
	if err != nil {
		logrus.WithError(err).Error("can not compile keymap")
		return err
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+compiler.FirmwareName)
	w.Header().Set("Content-Type", "application/octet-stream")
	written, err := io.Copy(w, firmware)
	if err != nil {
		logrus.WithError(err).WithField("writter size", written).Error("get error while writing response")
		return err
	}
	return nil
}
