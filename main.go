package main

import (
	"flag"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/tdycwym/edgex_admin/caller"
	"github.com/tdycwym/edgex_admin/config"
	"github.com/tdycwym/edgex_admin/logs"

	"github.com/tdycwym/edgex_admin/middleware/cors"
	"github.com/tdycwym/edgex_admin/middleware/session"
	"go.uber.org/zap"
)

func main() {

	var confFilePath string
	flag.StringVar(&confFilePath, "conf", "", "Specify local configuration file path")
	flag.Parse()

	config.LoadConfig(confFilePath)
	logs.InitLogs()
	caller.InitClient()

	gin.SetMode(config.Server.RunMode)

	r := gin.New()
	r.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(zap.L(), true))

	// 允许跨域访问
	r.Use(cors.CorsMiddleware())

	r.Use(session.EnableRedisSession())

	registerRouter(r)

	_ = r.Run(config.Server.Port)
}
