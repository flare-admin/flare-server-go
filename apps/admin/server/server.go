package server

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/flare-admin/flare-server-go/apps/admin/service"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/pkg/hredis"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/i18n"
	psb "github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/casbin"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/cors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/oplog"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/sql_injection"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support"
	storage_rest "github.com/flare-admin/flare-server-go/framework/support/storage/interfaces/rest"
	"github.com/hertz-contrib/gzip"
	"golang.org/x/text/language"
)

const baseUrl = "/api/admin"

func NewCasBinEnforcer(hc *hredis.RedisClient, pr psb.IPermissionsRepository) (*psb.Enforcer, error) {
	enforcer, err := psb.NewEnforcer(pr, hc, baseUrl)
	if err != nil {
		return nil, err
	}
	return enforcer, nil
}
func NewServer(
	config *configs.Bootstrap,
	tk token.IToken,
	oplDbWriter oplog.IDbOperationLogWrite,
	frameworkServer *support.Server,
	soc *service.SysCronService,
	// 文件存储服务
	fs *storage_rest.Service,

) *hserver.Serve {
	// 设置时区
	err := utils.SetTimeZone(config.Server.TimeZone)
	if err != nil {
		panic(err)
	}
	svr := hserver.NewServe(&hserver.ServerConfig{
		Port:               config.Server.Port,
		RateQPS:            config.Server.RateQPS,
		TracerPort:         config.Server.TracerPort,
		Name:               config.Server.Name,
		MaxRequestBodySize: config.Server.MaxRequestBodySize,
	}, hserver.WithTokenizer(tk), hserver.WithBaseUrl(baseUrl))
	registerMiddleware(config, svr.GetHertz(), oplDbWriter)
	svr.RegisterRouters(frameworkServer, fs)
	soc.Start()
	return svr
}
func registerMiddleware(con *configs.Bootstrap, server *server.Hertz, oplDbWriter oplog.IDbOperationLogWrite) {
	// Set up cross domain and flow limiting middleware
	server.Use(cors.Handler())
	//Use compression
	server.Use(gzip.Gzip(gzip.DefaultCompression))
	//internationalization
	if con.ConfPath != "" {
		ph := fmt.Sprintf("%s/localize", con.ConfPath)
		server.Use(i18n.Handler(ph, language.Chinese, language.Chinese, language.English, language.TraditionalChinese))
	}
	// server.Use(ratelimit.RateLimitMiddleware(10))
	// 防止sql注入
	server.Use(sql_injection.PreventSQLInjection())

	// 操作日志
	//initOpLog(con.Log)
	initDbOpLog(oplDbWriter)
}

func initDbOpLog(oplDbWriter oplog.IDbOperationLogWrite) {
	writer := oplog.NewDBWriter(oplDbWriter)
	oplog.Init(writer)
}
