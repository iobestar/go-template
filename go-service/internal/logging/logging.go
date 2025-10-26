package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/iobestar/go-template/go-service/internal"
	"github.com/natefinch/lumberjack"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

var (
	formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		},
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	}
)

func NewLogger(configs ...func(logger *logrus.Logger)) *logrus.Logger {
	logger := logrus.New()
	for _, config := range configs {
		config(logger)
	}
	return logger
}

func NewLevelWriter(level logrus.Level) io.WriteCloser {
	return logrus.StandardLogger().WriterLevel(level)
}

func Init(config Config) (internal.DisposeFunc, error) {
	logger := logrus.StandardLogger()
	logger.SetFormatter(formatter)
	logger.SetOutput(os.Stdout)

	if config.File.Enabled {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   config.File.Name,
			MaxSize:    config.File.MaxSize, // megabytes
			MaxAge:     config.File.MaxAge,  // days
			MaxBackups: config.File.MaxBackups,
		})
	}

	if l, err := logrus.ParseLevel(config.Level); nil != err {
		return nil, err
	} else {
		logrus.SetLevel(l)
	}

	logger.AddHook(newMetricsHook())
	return nil, nil
}

func newMetricsHook() logrus.Hook {
	counter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "logs_total",
			Help: "Number of logs by log level.",
		},
		[]string{"level"},
	)
	return &metricsHook{
		counter: counter,
	}
}

type metricsHook struct {
	counter *prometheus.CounterVec
}

func (h *metricsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *metricsHook) Fire(e *logrus.Entry) error {
	h.counter.WithLabelValues(e.Level.String()).Inc()
	return nil
}
