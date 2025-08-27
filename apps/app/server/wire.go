package server

import (
	"github.com/flare-admin/flare-server-go/apps/app/service"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/pkg/hredis"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/support"
	"github.com/flare-admin/flare-server-go/framework/support/storage"
	"github.com/flare-admin/flare-server-go/framework/support/storage/application"
	storage_rest "github.com/flare-admin/flare-server-go/framework/support/storage/interfaces/rest"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// 基础设施
	support.BaseProviderSet,
	sysevent.ProviderSet,

	NewServer,
	service.ProviderSet,

	NewRdbToken,
	storage.ProviderSet,
	NewFileService,
)

func NewRdbToken(config *configs.Bootstrap, hc *hredis.RedisClient) token.IToken {
	return token.NewRdbToken(hc.GetClient(), config.JWT.Issuer, config.JWT.SigningKey, config.JWT.ExpirationToken, config.JWT.ExpirationRefresh, true)
}

// NewFileService 文件上传服务
func NewFileService(ser *application.StorageService) *storage_rest.Service {
	return storage_rest.NewService("/api/admin/", ser)
}
