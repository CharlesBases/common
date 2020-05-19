package http_context

import (
	"sync"
	"testing"
)

func TestContext(t *testing.T) {
	waitgroup := sync.WaitGroup{}
	waitgroup.Add(4)

	context := NewContext()

	go func() {
		defer waitgroup.Done()

		context.Debug("debug", "debug")
		context.Debugf("debugf %s", "debug")
	}()

	go func() {
		defer waitgroup.Done()

		context.Info("info", "info")
		context.Infof("infof %s", "info")
	}()

	go func() {
		defer waitgroup.Done()

		context.Warn("warn", "warn")
		context.Warnf("warnf %s", "warn")
	}()

	go func() {
		defer waitgroup.Done()

		context.Error("error", "error")
		context.Errorf("errorf %s", "errorf")
	}()

	waitgroup.Wait()
	context.Flush()
}
