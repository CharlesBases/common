package memory

import (
	"fmt"
	"testing"
	"time"

	"github.com/CharlesBases/common/sync"
)

var mSync = NewSync(sync.WithTimeout(time.Second * 10))

func TestLock(t *testing.T) {
	mSync.Lock("test")

	go func() {
		if err := mSync.Lock("test", sync.WithLockTTL(time.Second*5)); err != nil {
			fmt.Println("1", err)
		} else {
			fmt.Println("lock 1")
		}
	}()

	go func() {
		if err := mSync.Lock("test", sync.WithLockTTL(time.Second*5)); err != nil {
			fmt.Println("2", err)
		} else {
			fmt.Println("lock 2")
		}
	}()

	go func() {
		for {
			select {
			case <-time.Tick(time.Second * 2):
				fmt.Println("unlock test")
				mSync.Unlock("test")
			}
		}
	}()

	<-time.Tick(time.Hour)
}

func TestUnlock(t *testing.T) {
	<-time.Tick(time.Hour)
}
