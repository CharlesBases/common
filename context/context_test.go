package context

import (
	"sync"
	"testing"
)

func TestContext(t *testing.T) {
	waitgroup := sync.WaitGroup{}
	waitgroup.Add(5)

	context := NewContext()

	go func() {
		defer waitgroup.Done()

		context.logger.Debug("debug", "debug")
		context.logger.Debugf("debugf %s", "debug")
	}()

	go func() {
		defer waitgroup.Done()

		context.logger.Info("info", "info")
		context.logger.Infof("infof %s", "info")
	}()

	go func() {
		defer waitgroup.Done()

		context.logger.Warn("warn", "warn")
		context.logger.Warnf("warnf %s", "warn")
	}()

	go func() {
		defer waitgroup.Done()

		context.logger.Error("error", "error")
		context.logger.Errorf("errorf %s", "warn")
	}()

	go func() {
		defer waitgroup.Done()

		context.logger.Critical("critical", "critical")
		context.logger.Criticalf("criticalf %s", "criticalf")
	}()

	waitgroup.Wait()
	context.Flush()
}
