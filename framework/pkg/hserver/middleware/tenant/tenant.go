package tenant

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
)

// IgnoreTenantHandler 租户处理
func IgnoreTenantHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ctx = actx.BuildIgnoreTenantCtx(ctx)
		c.Next(ctx)
	}
}
