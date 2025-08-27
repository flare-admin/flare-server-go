package support

import (
	"github.com/flare-admin/flare-server-go/framework/infrastructure"
	"github.com/flare-admin/flare-server-go/framework/support/base"
	"github.com/flare-admin/flare-server-go/framework/support/cache"
	"github.com/flare-admin/flare-server-go/framework/support/config_center"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary"
	"github.com/flare-admin/flare-server-go/framework/support/monitoring"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent"
	"github.com/flare-admin/flare-server-go/framework/support/systask"
	"github.com/flare-admin/flare-server-go/framework/support/template"
	"github.com/google/wire"
)

// ProviderSet 基础依赖
var ProviderSet = wire.NewSet(
	infrastructure.ProviderSet,
	base.ProviderSet,
	monitoring.ProviderSet,
	systask.ProviderSet,
	sysevent.ProviderSet,
	config_center.ProviderSet,
	cache.ProviderSet,
	dictionary.AdminProviderSet,
	template.ProviderSet,
	rule_engine.ProviderSet,
	NewServer,
)

// BaseProviderSet 基础依赖
var BaseProviderSet = wire.NewSet(
	infrastructure.ProviderSet,
	base.BaseProviderSet,
	config_center.BaseProviderSet,
	cache.BaseProviderSet,
	rule_engine.BaseProviderSet,
	dictionary.ProviderSet,
	template.BaseProviderSet,
)
