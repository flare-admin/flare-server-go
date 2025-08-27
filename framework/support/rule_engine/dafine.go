package rule_engine

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/interfaces/admin"
)

type RuleEngineServer struct {
	ts *admin.TemplateService
	rc *admin.CategoryService
	rs *admin.RuleService
}

func NewServer(
	ts *admin.TemplateService,
	rc *admin.CategoryService,
	rs *admin.RuleService,
) *RuleEngineServer {
	return &RuleEngineServer{
		rc: rc,
		rs: rs,
		ts: ts,
	}
}

func (s *RuleEngineServer) RegisterRouter(rg *route.RouterGroup, tk token.IToken) {
	s.rc.RegisterRouter(rg, tk)
	s.ts.RegisterRouter(rg, tk)
	s.rs.RegisterRouter(rg, tk)
}
