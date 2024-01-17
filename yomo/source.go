package yomo

import "log/slog"

// Source is responsible for sending data to yomo.
type Source interface {
	// Close will close the connection to YoMo-Zipper.
	// Close() error
	// Connect to YoMo-Zipper.
	Connect() error
	// Write the data to directed downstream.
	Write(tag uint32, data []byte) error
	// SetErrorHandler set the error handler function when server error occurs
	// SetErrorHandler(fn func(err error))
}

type source struct {
	name       string
	zipperAddr string
}

func NewSource(name string, zipperAddr string) Source {
	return &source{name: name, zipperAddr: zipperAddr}
}

func (s *source) Connect() error {
	slog.Info("source connect", "name", s.name, "zipperAddr", s.zipperAddr)
	return nil
}

func (s *source) Write(tag uint32, data []byte) error {
	slog.Info("source write", "name", s.name, "tag", tag, "data", string(data))
	return nil
}
