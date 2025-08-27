package dictionary

import (
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/biz"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/data"
	dictionaryinterfaces "github.com/flare-admin/flare-server-go/framework/support/dictionary/interfaces"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/translator"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	data.NewDictionaryRepo,
	biz.NewDictionaryUseCase,
	translator.NewTranslator,
)

var AdminProviderSet = wire.NewSet(
	ProviderSet,
	dictionaryinterfaces.NewDictionaryService,
)
