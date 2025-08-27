package main

import (
	"flag"
	_ "github.com/flare-admin/flare-server-go/apps/app/docs"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/pkg/constant"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/log"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/sql_injection"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
	"golang.org/x/exp/rand"
	"time"
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
	server *hserver.Serve
}

func newApp(server *hserver.Serve) *app {
	return &app{
		server: server,
	}
}

// @title go-server-template-app
// @version 1.0
// @description This is a demo using go-server-template-app.

// @contact.name go-server-template-app

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8848
// @BasePath /apps/app
// @schemes http
func main() {
	flag.Parse()
	//load config
	bc, err := configs.Load(flagConf, flagLocalize, env, flagLog)
	if err != nil {
		panic(err)
	}
	// 生成 6 位随机数
	rand.Seed(uint64(time.Now().UnixNano()))
	application, cleanup, err := wireApp(bc, bc.Data)
	if err != nil {
		panic(err)
	}
	hserver.InitLog(log.NewLogConf(
		log.WithOutPath(bc.Log.OutPath),
		log.WithLevel(bc.Log.Level),
		log.WithMaxSize(bc.Log.MaxSize),
		log.WithMaxAge(bc.Log.MaxAge),
		log.WithMaxBackups(bc.Log.MaxBackups),
		log.WithCompress(bc.Log.Compress),
		log.WithFilePrefix(bc.Log.FilePrefix),
	), constant.EnvMode(env))

	// 放置csrf攻击
	//store := cookie.NewStore([]byte("secret"))
	//server.Use(sessions.New("csrf-session", store))
	//server.Use(csrf.New())
	// 防止sql注入
	application.server.Use(sql_injection.PreventSQLInjection())
	url := swagger.URL("/swagger/doc.json") // The url pointing to API definition
	application.server.GetHertz().GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))
	defer cleanup()
	application.server.Run()
}
