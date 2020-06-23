package ilock

import (
	"context"
)

// 分布式锁

type Locker interface {
	// 返回true，表示成功lock，false已经上锁过，此次没有成功
	TryLock(ctx context.Context, lockId string, expireSecond int64) (lockSucceed bool, err error)
	// 获取基于lockId的锁
	// expireSecond秒后自动释放，为0则不自动释放;最多等待maxWaitSecond秒，必须大于0
	// 返回error为空，意味着已经获取到lock；否则为未获取到lock，或者出现其他error
	Lock(ctx context.Context, lockId string, expireSecond int64, maxWaitSecond uint32) <-chan error
	// 释放锁，idempotent
	Release(ctx context.Context, lockId string) (err error)
}
