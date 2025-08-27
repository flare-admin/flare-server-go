package goroutine

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"runtime/debug"
)

// Go 携程
func Go(x func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				hlog.Errorf("time:%s, err:%v, fatal%s", utils.GetTimeNow().Format("2006-01-02 15:04:05:06"), err, string(debug.Stack()))
			}
		}()
		x()
	}()
}

// SecureGo 安全的携程
func SecureGo(ctx context.Context, x func(context.Context)) {
	go func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				hlog.CtxErrorf(ctx, "time:%s, err:%v, fatal%s", utils.GetTimeNow().Format("2006-01-02 15:04:05:06"), err, string(debug.Stack()))
			}
		}()
		x(ctx)
	}(ctx)
}
