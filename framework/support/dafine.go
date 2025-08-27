package support

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/support/base"
	cache "github.com/flare-admin/flare-server-go/framework/support/cache/interfaces/rest"
	configcenter "github.com/flare-admin/flare-server-go/framework/support/config_center/interfaces/rest"
	monrest "github.com/flare-admin/flare-server-go/framework/support/monitoring/interfaces/rest"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine"
	syseventservice "github.com/flare-admin/flare-server-go/framework/support/sysevent/interfaces"
	systaskinterfaces "github.com/flare-admin/flare-server-go/framework/support/systask/interfaces"
	"github.com/flare-admin/flare-server-go/framework/support/template"
)

type Server struct {
	metrice  *monrest.MetricsController
	base     *base.BaseServer
	ccs      *configcenter.ConfigHandler
	cs       *cache.CacheHandler
	ts       *template.TempServer
	rs       *rule_engine.RuleEngineServer
	systask  *systaskinterfaces.TaskService
	sysevent *syseventservice.EventService
}

func NewServer(
	metrice *monrest.MetricsController,
	base *base.BaseServer,
	ccs *configcenter.ConfigHandler,
	cs *cache.CacheHandler,
	ts *template.TempServer,
	rs *rule_engine.RuleEngineServer,
	systask *systaskinterfaces.TaskService,
	sysevent *syseventservice.EventService,
) *Server {
	return &Server{
		metrice:  metrice,
		base:     base,
		ccs:      ccs,
		cs:       cs,
		ts:       ts,
		rs:       rs,
		systask:  systask,
		sysevent: sysevent,
	}
}

func (s *Server) RegisterRouter(rg *route.RouterGroup, tk token.IToken) {
	s.metrice.RegisterRouter(rg, tk)
	s.base.RegisterRouter(rg, tk)
	s.ccs.RegisterRouter(rg, tk)
	s.cs.RegisterRouter(rg, tk)
	s.ts.RegisterRouter(rg, tk)
	s.rs.RegisterRouter(rg, tk)

	s.systask.RegisterRouter(rg, tk)
	s.sysevent.RegisterRouter(rg, tk)
}
