package handlers

import (
	"context"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/dto"

	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/base/application/queries"
)

type DepartmentQueryHandler struct {
	queryService query.IDepartmentQueryService
}

func NewDepartmentQueryHandler(
	queryService query.IDepartmentQueryService,
) *DepartmentQueryHandler {
	return &DepartmentQueryHandler{
		queryService: queryService,
	}
}

// HandleList 处理部门列表查询
func (h *DepartmentQueryHandler) HandleList(ctx context.Context, req *queries.ListDepartmentsQuery) (*models.PageRes[dto.DepartmentDto], herrors.Herr) {
	if validate := req.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Query validation error: %s", validate)
		return nil, validate
	}

	// 构建查询条件
	qb := db_query.NewQueryBuilder()
	if req.Name != "" {
		qb.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	if req.Code != "" {
		qb.Where("code", db_query.Like, "%"+req.Code+"%")
	}
	if req.Status != nil {
		qb.Where("status", db_query.Eq, *req.Status)
	}
	if req.ParentID != "" {
		qb.Where("parent_id", db_query.Eq, req.ParentID)
	}
	qb.OrderBy("sequence", false)
	qb.WithPage(&req.Page)

	// 查询总数
	total, err := h.queryService.CountDepartments(ctx, qb)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to count departments: %s", err)
		return nil, herrors.QueryFail(err)
	}

	// 查询列表数据
	depts, err := h.queryService.FindDepartments(ctx, qb)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to list departments: %s", err)
		return nil, herrors.QueryFail(err)
	}

	// 返回分页结果
	return &models.PageRes[dto.DepartmentDto]{
		Total: total,
		List:  depts,
	}, nil
}

// HandleGet 处理获取部门查询
func (h *DepartmentQueryHandler) HandleGet(ctx context.Context, query *queries.GetDepartmentQuery) (*dto.DepartmentDto, herrors.Herr) {
	if validate := query.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Query validation error: %s", validate)
		return nil, validate
	}

	// 查询部门
	dept, err := h.queryService.GetDepartment(ctx, query.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get department: %s", err)
		return nil, herrors.QueryFail(err)
	}
	if dept == nil {
		return nil, herrors.QueryFail(fmt.Errorf("department not found: %s", query.ID))
	}

	return dept, nil
}

// HandleGetTree 处理获取部门树查询
func (h *DepartmentQueryHandler) HandleGetTree(ctx context.Context, query *queries.GetDepartmentTreeQuery) ([]*dto.DepartmentTreeDto, herrors.Herr) {
	if validate := query.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Query validation error: %s", validate)
		return nil, validate
	}

	// 获取部门树
	tree, err := h.queryService.GetDepartmentTree(ctx, query.ParentID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get department tree: %s", err)
		return nil, herrors.QueryFail(err)
	}

	return tree, nil
}

// HandleGetUserDepartments 处理获取用户部门查询
func (h *DepartmentQueryHandler) HandleGetUserDepartments(ctx context.Context, query *queries.GetUserDepartmentsQuery) ([]*dto.DepartmentDto, herrors.Herr) {
	if validate := query.Validate(); herrors.HaveError(validate) {
		hlog.CtxErrorf(ctx, "Query validation error: %s", validate)
		return nil, validate
	}

	// 获取用户部门
	depts, err := h.queryService.GetUserDepartments(ctx, query.UserID)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get user departments: %s", err)
		return nil, herrors.QueryFail(err)
	}

	return depts, nil
}

// HandleGetUsers 处理获取部门用户
func (h *DepartmentQueryHandler) HandleGetUsers(ctx context.Context, req *queries.GetDepartmentUsersQuery) (*models.PageRes[dto.UserDto], herrors.Herr) {
	// 1. 构建查询条件
	qb := db_query.NewQueryBuilder()
	if req.Username != "" {
		qb.Where("username", db_query.Like, "%"+req.Username+"%")
	}
	if req.Name != "" {
		qb.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	qb.WithPage(&req.Page)
	qb.OrderBy("created_at", true)

	// 2. 查询总数
	total, err := h.queryService.CountDepartmentUsers(ctx, req.DeptID, "", qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 3. 查询用户列表
	users, err := h.queryService.GetDepartmentUsers(ctx, req.DeptID, "", qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 4. 返回分页结果
	return &models.PageRes[dto.UserDto]{
		List:  users,
		Total: total,
	}, nil
}

// HandleGetUnassignedUsers 处理获取未分配部门的用户查询
func (h *DepartmentQueryHandler) HandleGetUnassignedUsers(ctx context.Context, req *queries.GetUnassignedUsersQuery) (*models.PageRes[dto.UserDto], herrors.Herr) {
	// 构建查询条件
	qb := db_query.NewQueryBuilder()
	if req.Username != "" {
		qb.Where("username", db_query.Like, "%"+req.Username+"%")
	}
	if req.Name != "" {
		qb.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	qb.Where("status", db_query.Eq, 1)
	qb.WithPage(&req.Page)
	qb.OrderBy("created_at", true)

	// 查询总数
	total, err := h.queryService.CountUnassignedUsers(ctx, qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 查询用户列表
	users, err := h.queryService.GetUnassignedUsers(ctx, qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	return &models.PageRes[dto.UserDto]{
		List:  users,
		Total: total,
	}, nil
}
