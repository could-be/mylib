package userctx

import "context"

type keyType int

const (
	uidkey          keyType = 1
	platformKey     keyType = 2
	deviceKey       keyType = 3
	deviceSourceKey keyType = 4
	InvalidUid      uint64  = 0
)

func WithContext(parent context.Context, uid uint64) context.Context {
	return context.WithValue(parent, uidkey, uid)
}

/* 获取buyer_id*/
func FromContext(ctx context.Context) (uid uint64, ok bool) {
	uid, ok = ctx.Value(uidkey).(uint64)
	return
}

func WithContextPlatform(parent context.Context, plt int32) context.Context {
	return context.WithValue(parent, platformKey, plt)
}
func FromContextPlatform(ctx context.Context) (int32, bool) {
	plt, ok := ctx.Value(platformKey).(int32)
	return plt, ok
}

func WithContextDevice(parent context.Context, device string) context.Context {
	return context.WithValue(parent, deviceKey, device)
}
func FromContextDevice(ctx context.Context) (string, bool) {
	device, ok := ctx.Value(deviceKey).(string)
	return device, ok
}

func WithContextDeviceSource(parent context.Context, deviceSource string) context.Context {
	return context.WithValue(parent, deviceSourceKey, deviceSource)
}
func FromContextDeviceSource(ctx context.Context) (string, bool) {
	deviceSource, ok := ctx.Value(deviceSourceKey).(string)
	return deviceSource, ok
}
