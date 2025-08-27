package actx

import (
	"context"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/constant"
	"golang.org/x/exp/slices"
	"strings"

	"github.com/flare-admin/flare-server-go/framework/pkg/token"
)

const (
	keyAccessToken = "access_token"
	KeyUserId      = "userId"
	KeyUsername    = "username"
	KeyPlatform    = "platform"
	KeyToken       = "token"
	KeyRole        = "role"
	KeyTenantId    = "tenant_id"
	DeviceId       = "deviceId"
	DeviceName     = "deviceName"
	IpAddress      = "ipAddress"
	UserAgent      = "UserAgent"
	IgnoreTenantId = "ignore_tenant_Id"
	KeyDeptId      = "deptId"
)

func WithUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, KeyUserId, userId)
}

func GetUserId(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyUserId))
}
func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, KeyUsername, username)
}
func GetUsername(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyUsername))
}
func WithPlatform(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, KeyPlatform, platform)
}

func GetPlatform(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyPlatform))
}
func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, KeyToken, token)
}

func GetDeptId(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyToken))
}

func WithDeptId(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, KeyDeptId, token)
}

func GetToken(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyDeptId))
}
func WithRole(ctx context.Context, role []string) context.Context {
	return context.WithValue(ctx, KeyRole, strings.Join(role, ","))
}

func GetRoles(ctx context.Context) []string {
	value := fmt.Sprintf("%v", ctx.Value(KeyRole))
	if value == "" {
		return []string{}
	}
	return strings.Split(value, ",")
}
func WithTenantId(ctx context.Context, tenantId string) context.Context {
	return context.WithValue(ctx, KeyTenantId, tenantId)
}

func GetTenantId(ctx context.Context) string {
	tenId := fmt.Sprintf("%v", ctx.Value(KeyTenantId))
	if tenId == IgnoreTenantId {
		return ""
	}
	return tenId
}

func WithDeviceId(ctx context.Context, deviceId string) context.Context {
	return context.WithValue(ctx, DeviceId, deviceId)
}

func GetDeviceId(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(DeviceId))
}
func WithDeviceName(ctx context.Context, deviceName string) context.Context {
	return context.WithValue(ctx, DeviceName, deviceName)
}

func GetDeviceName(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(DeviceName))
}

func WithIpAddress(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, IpAddress, addr)
}

func GetIpAddress(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(IpAddress))
}

func WithUserAgent(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, UserAgent, addr)
}

func GetUserAgent(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(UserAgent))
}

func WithIgnoreTenantId(ctx context.Context) context.Context {
	return context.WithValue(ctx, IgnoreTenantId, IgnoreTenantId)
}
func GetIgnoreTenantId(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(IgnoreTenantId))
}
func IsIgnoreTenantId(ctx context.Context) bool {
	return GetIgnoreTenantId(ctx) == IgnoreTenantId
}

// BuildIgnoreTenantCtx 构建忽略租户的ctx
func BuildIgnoreTenantCtx(ctx context.Context) context.Context {
	return WithIgnoreTenantId(ctx)
}

func Store(ctx context.Context, accessToken token.AccessToken) context.Context {
	ctx = WithUserId(ctx, accessToken.UserId)
	ctx = WithPlatform(ctx, accessToken.Platform)
	ctx = WithToken(ctx, accessToken.AccessToken)
	ctx = WithRole(ctx, accessToken.Roles)
	ctx = WithTenantId(ctx, accessToken.TenantId)
	ctx = WithUsername(ctx, accessToken.UserName)
	return ctx
}
func IsSuperAdmin(ctx context.Context) bool {
	roles := GetRoles(ctx)
	if len(roles) == 0 {
		return false
	}
	return slices.Contains(roles, constant.RoleSuperAdmin)
}
func IsAdmin(ctx context.Context) bool {
	role := GetRoles(ctx)
	if len(role) == 0 {
		return false
	}
	return true
	//if slices.Contains(role, constant.RoleSuperAdmin) {
	//	return true
	//}
	//if slices.Contains(role, constant.RoleAdmin) {
	//	return true
	//}
	//if slices.Contains(role, constant.RoleAgent) {
	//	return true
	//}
	//return false
}
