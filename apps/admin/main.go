package main

import (
	"flag"
	"time"

	_ "github.com/flare-admin/flare-server-go/apps/admin/docs"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/pkg/constant"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/log"
	"github.com/flare-admin/flare-server-go/framework/pkg/mqevent/manager"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
	"golang.org/x/exp/rand"

	// 确保 swag 扫描到 rule_engine admin 接口
	_ "github.com/flare-admin/flare-server-go/framework/support/rule_engine/interfaces/admin"
)

var (
	// flagconf is the config flag.
	flagConf string
	// flagconf is the config flag.
	env string
	// flagLocalize is the Localize flag.
	flagLocalize string
	// flagLog is the login flag.
	flagLog string
)

func init() {
	flag.StringVar(&flagConf, "conf", "../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&env, "env", "dev", "Operating environment, eg: -env dev")
	flag.StringVar(&flagLocalize, "lc", "", "localize config path, eg: -lc localize")
	flag.StringVar(&flagLog, "log", "", "Operating environment, eg: -log app.log")
}

type app struct {
	server      *hserver.Serve
	eventServer manager.EventManager
}

func newApp(server *hserver.Serve, eventServer manager.EventManager) *app {
	return &app{
		server:      server,
		eventServer: eventServer,
	}
}

// @title go-server-template-admin
// @version 1.0
// @description This is a demo using go-server-template-admin.

// @contact.name go-server-template-admin

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8858
// @BasePath /
// @schemes http
func main() {
	flag.Parse()
	//load config
	bc, err := configs.Load(flagConf, flagLocalize, env, flagLog)
	if err != nil {
		panic(err)
	}
	bc.ConfPath = flagConf
	hserver.InitLog(log.NewLogConf(
		log.WithOutPath(bc.Log.OutPath),
		log.WithLevel(bc.Log.Level),
		log.WithMaxSize(bc.Log.MaxSize),
		log.WithMaxAge(bc.Log.MaxAge),
		log.WithMaxBackups(bc.Log.MaxBackups),
		log.WithCompress(bc.Log.Compress),
		log.WithFilePrefix(bc.Log.FilePrefix),
	), constant.EnvMode(env))
	// 生成 6 位随机数
	rand.Seed(uint64(time.Now().UnixNano()))
	application, cleanup, err := wireApp(bc, bc.Data)
	if err != nil {
		panic(err)
	}
	url := swagger.URL("/swagger/doc.json") // The url pointing to API definition
	application.server.GetHertz().GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))
	defer func() {
		cleanup()
		application.eventServer.Close()
	}()
	application.eventServer.Start()
	application.server.Run()

}
