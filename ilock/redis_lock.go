package ilock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gopkg.in/redis.v5"
	"ptapp.cn/util/iredis"
)

const (
	redisPrefix   = "ilock:"
	getset_expire = `
local val, err = redis.pcall('GETSET', KEYS[1], '123')
if err then
    return err
end
redis.call('EXPIRE', KEYS[1], ARGV[1])
if val then
    return {val}
end
return nil
`
)

type redisLock struct {
	cli          *iredis.Client
	getsetexpCmd *redis.Script
}

type RedisLockConf struct {
	Redis iredis.RedisCfg `json:"redis"`
}

func NewRedisLock(cfg *RedisLockConf) (Locker, error) {
	cli, err := iredis.NewClient(cfg.Redis)
	if err != nil {
		return nil, err
	}
	return &redisLock{
		cli:          cli,
		getsetexpCmd: redis.NewScript(getset_expire),
	}, nil
}

func (r *redisLock) TryLock(ctx context.Context, id string, expireSecond int64) (lockSucceed bool, err error) {

	err = r.getsetexpCmd.Run(r.cli, []string{redisPrefix + id}, expireSecond).Err()
	if err != nil {
		if err != redis.Nil {
			return false, err
		} else {
			return true, nil
		}

	}
	return false, nil
}

var waitInterval time.Duration = time.Millisecond * 50 // 每次等待50ms
func (r *redisLock) Lock(ctx context.Context, lockId string, expireSecond int64, maxWaitSecond uint32) <-chan error {
	ret := make(chan error, 2)
	go func() {
		defer close(ret)
		defer func() {
			if e := recover(); e != nil {
				ret <- fmt.Errorf("panic:%v", e)
			}
		}()

		if maxWaitSecond == 0 {
			ret <- errors.New("zero maxWaitSecond")
			return
		}

		locked := false
		var err error
		t := time.NewTimer(time.Duration(maxWaitSecond) * time.Second)
		for {
			locked, err = r.TryLock(ctx, lockId, expireSecond)
			if err != nil {
				ret <- err
				return
			}
			if locked {
				break
			}
			select {
			case <-t.C:
				ret <- errors.New("exceeds maxWaitSecond")
				return
			case <-time.After(waitInterval):
			}
		}
		ret <- nil
	}()
	return ret
}

func (r *redisLock) Release(ctx context.Context, lockId string) (err error) {
	err = r.cli.Del(redisPrefix + lockId).Err()
	return
}
