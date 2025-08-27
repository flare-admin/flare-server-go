package template

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/support/template/interfaces/admin"
)

type TempServer struct {
	cs *admin.CategoryService
	ts *admin.TemplateService
}

func NewTempServer(
	cs *admin.CategoryService,
	ts *admin.TemplateService,

) *TempServer {
	return &TempServer{
		cs: cs,
		ts: ts,
	}
}

func (s *TempServer) RegisterRouter(rg *route.RouterGroup, tk token.IToken) {
	s.cs.RegisterRouter(rg, tk)
	s.ts.RegisterRouter(rg, tk)
}
