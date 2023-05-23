package log

import (
	"time"

	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/sirupsen/logrus"
)

const defaultTimestampFormat = time.RFC3339

type BloomLogFormatter struct {
	prefixed.TextFormatter
	BloomDisbleTs bool
}

func (f *BloomLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	f.DisableTimestamp = true
	b, err := f.TextFormatter.Format(entry)
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if !f.BloomDisbleTs {
		logTime := entry.Time.Format(timestampFormat) + " "
		b = append([]byte(logTime), b...)
	}
	return b, err
}
