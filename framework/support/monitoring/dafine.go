package monitoring

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	monrest "github.com/flare-admin/flare-server-go/framework/support/monitoring/interfaces/rest"
)

type Server struct {
	metrice *monrest.MetricsController
}

func NewServer(
	metrice *monrest.MetricsController,
) *Server {
	return &Server{
		metrice: metrice,
	}
}

func (s *Server) RegisterRouter(rg *route.RouterGroup, tk token.IToken) {
	s.metrice.RegisterRouter(rg, tk)
}
