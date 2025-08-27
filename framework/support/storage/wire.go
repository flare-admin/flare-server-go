package storage

import (
	"github.com/flare-admin/flare-server-go/framework/support/storage/application"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain"
	"github.com/flare-admin/flare-server-go/framework/support/storage/infrastructure"
	"github.com/flare-admin/flare-server-go/framework/support/storage/infrastructure/persistence/data"
	"github.com/flare-admin/flare-server-go/framework/support/storage/infrastructure/persistence/repository"
	"github.com/flare-admin/flare-server-go/framework/support/storage/interfaces/api"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	data.NewFileRepository,
	repository.NewStorageRepo,
	infrastructure.NewStorageFactory,
	infrastructure.NewStorageAdapter,
	domain.NewStorageService,
	application.NewStorageService,
	storage_api.NewStorageApi,
)
