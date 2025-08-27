package biz

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
)

func getCtx(ctx context.Context) context.Context {
	return actx.BuildIgnoreTenantCtx(ctx)
}
