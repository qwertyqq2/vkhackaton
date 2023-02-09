package syncbc

type SyncBcMutex struct {
	ch chan struct{}
}

func NewSyncBc() *SyncBcMutex {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return &SyncBcMutex{
		ch: ch,
	}
}

func (sm *SyncBcMutex) TryLock() bool {
	_, ok := <-sm.ch
	return ok
}

func (sm *SyncBcMutex) Lock() {
	_, ok := <-sm.ch
	if !ok {
		panic("already locked")
	}
}

func (sm *SyncBcMutex) Unlock() {
	select {
	case sm.ch <- struct{}{}:

	default:
		panic("alreay unlocked")
	}
}

func (sm *SyncBcMutex) Locking() {
	for {
		if ok := sm.TryLock(); ok {
			return
		}
	}
}
