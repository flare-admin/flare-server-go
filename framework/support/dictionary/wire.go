package dictionary

import (
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/biz"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/data"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/translator"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	data.NewDictionaryRepo,
	biz.NewDictionaryUseCase,
	translator.NewTranslator,
)
