package redis

import (
	"log"
	"testing"
	"time"

	"github.com/CharlesBases/common/sync"
)

var rSync sync.Sync

func init() {
	rSync = NewStore(sync.WithBlocked())
}

// TestLock
func TestLock(t *testing.T) {
	if err := rSync.Lock("1", sync.WithLockTTL(time.Second*3)); err != nil {
		log.Fatal(err)
	}
}

// TestUnlock
func TestUnlock(t *testing.T) {
	rSync.Unlock("1")
}
