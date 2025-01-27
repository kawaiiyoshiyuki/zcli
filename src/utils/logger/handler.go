package logger

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Config struct {
	FilePath string
}

type Handler struct {
	config Config

	logrus *logrus.Logger
}

func New(config Config) (*Handler, error) {

	dir := filepath.Dir(config.FilePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	l := logrus.New()
	l.Out = ioutil.Discard
	l.Level = logrus.DebugLevel
	l.AddHook(&StdoutHook{levels: []logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel}})
	l.AddHook(&VarLogHook{
		path:   config.FilePath,
		levels: []logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel},
	})

	return &Handler{
		config: config,

		logrus: l,
	}, nil
}

func (h *Handler) Info(a ...interface{}) {
	h.logrus.Info(a...)
}

func (h *Handler) Infof(format string, a ...interface{}) {
	h.logrus.Infof(format, a...)
}

func (h *Handler) Warning(a ...interface{}) {
	h.logrus.Warning(a...)
}

func (h *Handler) Warningf(format string, a ...interface{}) {
	h.logrus.Warningf(format, a...)
}

func (h *Handler) Error(a ...interface{}) {
	h.logrus.Error(a...)
}

func (h *Handler) Errorf(format string, a ...interface{}) {
	h.logrus.Errorf(format, a...)
}

func (h *Handler) Debug(a ...interface{}) {
	h.logrus.Debug(a...)
}

func (h *Handler) Debugf(format string, a ...interface{}) {
	h.logrus.Debugf(format, a...)
}
