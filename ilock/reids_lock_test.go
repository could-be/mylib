package ilock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"ptapp.cn/util/iredis"
)

func TestRedisTryLock(t *testing.T) {
	ast := assert.New(t)
	ctx := context.Background()

	clearRedis()
	locker, err := NewRedisLock(&RedisLockConf{
		Redis: iredis.RedisCfg{
			Host: "localhost:6379",
			DB:   10,
		},
	})
	ast.Nil(err)

	ret, err := locker.TryLock(ctx, "haha", 1)
	ast.Nil(err)
	ast.True(ret)

	ret, err = locker.TryLock(ctx, "haha", 1)
	ast.Nil(err)
	ast.False(ret)
	time.Sleep(time.Second)
	ret, err = locker.TryLock(ctx, "haha", 1)
	ast.Nil(err)
	ast.True(ret)
}

func clearRedis() {
	cli, _ := iredis.NewClient(iredis.RedisCfg{
		Host: "localhost:6379",
		DB:   10,
	})
	cli.FlushDb()

}
