package syncbc

import (
	"fmt"
	"testing"
	"time"
)

func maingor(sm *SyncBcMutex) {
	for _, i := range []int{1, 2, 3, 4} {
		sm.Unlock()
		time.Sleep(1 * time.Second)
		sm.Lock()
		fmt.Println("Yes", i)
	}

}

func TestMutex(t *testing.T) {
	sm := NewSyncBc()
	fmt.Println(sm.TryLock())
	maingor(sm)
}
