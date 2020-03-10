package context

import (
	"fmt"

	"github.com/CharlesBases/common/algo"
	"github.com/CharlesBases/common/log"
)

type logger struct {
	trace string
}

func newlogger() *logger {
	return &logger{trace: fmt.Sprintf("[traceID: %s]", algo.GetTraceID())}
}

func (logger *logger) Debug(value ...interface{}) {
	log.Debug(append([]interface{}{logger.trace}, value...)...)
}

func (logger *logger) Debugf(format string, value ...interface{}) {
	log.Debugf(fmt.Sprintf("%s %s", logger.trace, format), value...)
}

func (logger *logger) Info(value ...interface{}) {
	log.Info(append([]interface{}{logger.trace}, value...)...)
}

func (logger *logger) Infof(format string, value ...interface{}) {
	log.Infof(fmt.Sprintf("%s %s", logger.trace, format), value...)
}

func (logger *logger) Warn(value ...interface{}) {
	log.Warn(append([]interface{}{logger.trace}, value...)...)
}

func (logger *logger) Warnf(format string, value ...interface{}) {
	log.Warnf(fmt.Sprintf("%s %s", logger.trace, format), value...)
}

func (logger *logger) Error(value ...interface{}) {
	log.Error(append([]interface{}{logger.trace}, value...)...)
}

func (logger *logger) Errorf(format string, value ...interface{}) {
	log.Errorf(fmt.Sprintf("%s %s", logger.trace, format), value...)
}

func (logger *logger) Critical(value ...interface{}) {
	log.Critical(append([]interface{}{logger.trace}, value...)...)
}

func (logger *logger) Criticalf(format string, value ...interface{}) {
	log.Criticalf(fmt.Sprintf("%s %s", logger.trace, format), value...)
}

func (logger *logger) flush() {
	log.Flush()
}
