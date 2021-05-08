package main

import (
	"flag"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/tydcwym/edgex_admin/caller"
	"github.com/tydcwym/edgex_admin/config"
	"github.com/tydcwym/edgex_admin/logs"

	"github.com/tydcwym/edgex_admin/middleware/cors"
	"github.com/tydcwym/edgex_admin/middleware/session"
	"go.uber.org/zap"

	"github.com/tydcwym/edgex_admin/handlers/user"
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
