package base

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/handlers"
	baserest "github.com/flare-admin/flare-server-go/framework/support/base/interfaces/rest"
)

type BaseServer struct {
	rc           *baserest.SysRoleController
	uc           *baserest.SysUserController
	ts           *baserest.SysTenantController
	ps           *baserest.SysPermissionsController
	as           *baserest.AuthController
	lls          *baserest.LoginLogController
	ols          *baserest.OperationLogController
	des          *baserest.DepartmentController
	dps          *baserest.DataPermissionController
	handlerEvent *handlers.HandlerEvent
}

func NewBaseServer(
	rc *baserest.SysRoleController,
	uc *baserest.SysUserController,
	ts *baserest.SysTenantController,
	ps *baserest.SysPermissionsController,
	as *baserest.AuthController,
	lls *baserest.LoginLogController,
	ols *baserest.OperationLogController,
	des *baserest.DepartmentController,
	dps *baserest.DataPermissionController,
	handlerEvent *handlers.HandlerEvent,
) *BaseServer {
	return &BaseServer{
		rc:           rc,
		uc:           uc,
		ts:           ts,
		ps:           ps,
		as:           as,
		lls:          lls,
		ols:          ols,
		des:          des,
		dps:          dps,
		handlerEvent: handlerEvent,
	}
}

func (s *BaseServer) RegisterRouter(rg *route.RouterGroup, tk token.IToken) {
	s.rc.RegisterRouter(rg, tk)
	s.uc.RegisterRouter(rg, tk)
	s.ts.RegisterRouter(rg, tk)
	s.ps.RegisterRouter(rg, tk)
	s.as.RegisterRouter(rg, tk)
	s.lls.RegisterRouter(rg, tk)
	s.ols.RegisterRouter(rg, tk)
	s.des.RegisterRouter(rg, tk)
	s.dps.RegisterRouter(rg, tk)
	s.handlerEvent.Register()
}
