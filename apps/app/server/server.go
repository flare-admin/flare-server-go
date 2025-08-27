package server

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/pkg/hredis"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/i18n"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/cors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/sql_injection"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/hertz-contrib/gzip"
	"golang.org/x/text/language"
)

const baseUrl = "/api/app"

func NewServer(
	config *configs.Bootstrap,
	hc *hredis.RedisClient,
) *hserver.Serve {
	// 设置时区
	err := utils.SetTimeZone(config.Server.TimeZone)
	if err != nil {
		panic(err)
	}

	tk := token.NewRdbToken(hc.GetClient(), config.JWT.Issuer, config.JWT.SigningKey, config.JWT.ExpirationToken, config.JWT.ExpirationRefresh, true)
	svr := hserver.NewServe(&hserver.ServerConfig{
		Port:               config.Server.Port,
		RateQPS:            config.Server.RateQPS,
		TracerPort:         config.Server.TracerPort,
		Name:               config.Server.Name,
		MaxRequestBodySize: config.Server.MaxRequestBodySize,
	}, hserver.WithTokenizer(tk), hserver.WithBaseUrl(baseUrl))
	registerMiddleware(config, svr.GetHertz())
	svr.RegisterRouters()
	return svr
}
func registerMiddleware(con *configs.Bootstrap, server *server.Hertz) {
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
}
