// Copyright 2020-2024 NGR Softlab
package emailer

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	logger = &logrus.Logger{
		Out:   io.MultiWriter(os.Stderr),
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05",
			ForceColors:            true,
			DisableLevelTruncation: true,
		},
		ReportCaller: true,
	}
)
