package jwt

import (
	"errors"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/flare-admin/flare-server-go/framework/pkg/constant"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
)

var (
	ErrNotFound = errors.New("not found from context")
)

// ParseToken 解析TOKEN
func ParseToken(c *app.RequestContext) (*token.AccessToken, error) {
	val, ok := c.Get(constant.KeyAccessToken)
	if !ok {
		return nil, ErrNotFound
	}

	accessToken, ok := val.(token.AccessToken)
	if !ok {
		return nil, fmt.Errorf("parse error")
	}

	return &accessToken, nil
}
