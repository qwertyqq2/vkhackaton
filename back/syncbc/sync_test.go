package syncbc

import (
	"fmt"
	"sync"
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

func Print() {
	fmt.Println(5)
}

func TestMutex(t *testing.T) {
	var once sync.Once
	once.Do(Print)
}
