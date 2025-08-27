package storage_rest

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/constant"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/storage/application"
	"mime"
	"path/filepath"
	"strings"
)

var (
	UploadFileErr        = herrors.New(constant.StatusServerError, "UploadFileErr", "upload File Err")
	UploadFileIsNotEmpty = herrors.New(constant.StatusServerError, "UploadFileIsNotEmpty", "Upload file is not empty")
	ObjectPathError      = herrors.New(constant.StatusServerError, "ObjectPathError", "路径错误")
)

// Service 文件上传
type Service struct {
	BaseApi string //基础地址
	ser     *application.StorageService
}

// NewService 文件上传
func NewService(baseApi string, ser *application.StorageService) *Service {
	return &Service{BaseApi: baseApi, ser: ser}
}

// RegisterRouter  路由注册
func (s *Service) RegisterRouter(r *route.RouterGroup, t token.IToken) {
	g := r.Group("file")
	{
		g.POST("singleFile", s.singleFile)
		g.POST("multiFile", s.multiFile)
		g.GET("object/:name", s.GetObject) //
	}
}

// singleFile 文件上传
// @Summary 文件上传
// @Description 文件上传
// @Tags 文件
// @ID singleFile
// @Accept application/json
// @Produce application/json
// @Success 200 {object} base_info.Success{data=UploadRes} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/app/file/singleFile [post]
func (s *Service) singleFile(ctx context.Context, c *app.RequestContext) {
	file, err := c.FormFile("file")
	if err != nil || file == nil {
		hlog.Errorf("get file err:%v", err)
		hserver.ResponseFailureErr(ctx, c, UploadFileIsNotEmpty)
		return
	}
	singleFile, err := s.ser.SingleFile(ctx, file)
	if err != nil {
		hlog.Errorf("upload file err:%v", err)
		hserver.ResponseFailureErr(ctx, c, UploadFileErr)
		return
	}
	hserver.ResponseSuccess(ctx, c, map[string]string{"url": singleFile})
}

// multiFile 多文件上传
// @Summary 多文件上传
// @Description 文件上多文件上传传
// @Tags 文件
// @ID multiFile
// @Accept application/json
// @Produce application/json
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router  /api/app/file/multiFile [post]
func (s *Service) multiFile(ctx context.Context, c *app.RequestContext) {
	form, err := c.MultipartForm()
	if err != nil {
		hlog.Errorf("upload file err:%v", err)
		hserver.ResponseFailureErr(ctx, c, UploadFileIsNotEmpty)
		return
	}
	files := form.File["file"]
	if len(files) == 0 {
		hlog.Errorf("file num zero ")
		hserver.ResponseFailureErr(ctx, c, UploadFileIsNotEmpty)
		return
	}
	urls, err := s.ser.MultiFile(ctx, files)
	if err != nil {
		hlog.Errorf("upload file err:%v", err)
		hserver.ResponseFailureErr(ctx, c, UploadFileErr)
		return
	}
	hserver.ResponseSuccess(ctx, c, map[string][]string{"urls": urls})
}

// GetObject 获取文件内容
// @Summary 获取文件内容
// @Description 获取文件内容，支持图片、HTML、文本等直接渲染，其他文件触发下载
// @Tags 文件
// @ID GetObject
// @Accept application/json
// @Produce application/json
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router  /api/app/file/object/:name [get]
func (s *Service) GetObject(ctx context.Context, c *app.RequestContext) {
	req := ObjectUrl{}
	if err := c.BindAndValidate(&req); err != nil {
		hserver.ResponseFailureErr(ctx, c, ObjectPathError)
		return
	}
	objectName := req.Name

	// 获取文件内容
	content, err := s.ser.ReadFileToContent(ctx, objectName)
	if err != nil {
		hserver.ResponseFailureErr(ctx, c, ObjectPathError)
		return
	}

	// 获取文件扩展名和 MIME 类型
	ext := strings.ToLower(filepath.Ext(objectName))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 设置通用响应头
	c.Response.Header.Set("Content-Type", contentType)

	// 判断是否是可以直接展示的文件类型
	if utils.IsViewableFile(ext) {
		// 直接在浏览器中渲染
		c.Response.Header.Set("Content-Disposition", "inline")
		c.Response.SetStatusCode(consts.StatusOK)
		c.Response.SetBody([]byte(content))
		return
	}

	// 其他文件触发下载
	c.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, objectName))
	c.Response.SetStatusCode(consts.StatusOK)
	c.Response.SetBody([]byte(content))
}

// ObjectUrl 参数绑定需要配合特定的 go tag 使用
type ObjectUrl struct {
	Name string `query:"name" path:"name"`
}
