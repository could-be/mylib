package util

import "time"

// 重试接口
func Retry(interval time.Duration,
	handle func() error) (err error) {

	for i := 0; i < 3; i++ {
		// 重试时间间隔
		<-time.After(interval)
		if err = handle(); err == nil {
			return
		}
	}

	return
}
