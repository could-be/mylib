package ilock

type Locker interface {
	Lock()
	TryLock()
}
