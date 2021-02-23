package cron

import (
	robfig "github.com/robfig/cron/v3"

	"github.com/CharlesBases/common/log"
)

var c = robfig.New()

// AddFunc .
func AddFunc(name string, spec string, cmd func()) {
	log.Debug("add task ", name, spec)

	c.AddFunc(spec, cmd)
}
