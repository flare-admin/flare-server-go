package base_api

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/dto"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query"
)

type ISysUserApi interface {
	// Get 获取用户信息
	Get(ctx context.Context, id string) (*dto.UserDto, error)
	// GetByInvitationCode 通过邀请码获取用户信息
	GetByInvitationCode(ctx context.Context, inviteCode string) (*dto.UserDto, error)
}

type SysUserApi struct {
	qr query.IUserQueryService
}

func NewSysUserApi(qr query.IUserQueryService) ISysUserApi {
	return &SysUserApi{qr: qr}
}

func (s SysUserApi) Get(ctx context.Context, id string) (*dto.UserDto, error) {
	return s.qr.GetUser(ctx, id)
}

func (s SysUserApi) GetByInvitationCode(ctx context.Context, inviteCode string) (*dto.UserDto, error) {
	return s.qr.GetByInvitationCode(ctx, inviteCode)
}
