package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tydcwym/edgex_admin/handlers"
	"github.com/tydcwym/edgex_admin/handlers/edgex"
	"github.com/tydcwym/edgex_admin/handlers/user"
	"github.com/tydcwym/edgex_admin/middleware/session"
	"github.com/tydcwym/edgex_admin/resp"
)

func registerRouter(r *gin.Engine) {
	r.GET("/ping", handlers.Ping)
	// your code

	edgexRouter := r.Group("/edgex_admin/edgex", session.AuthSessionMiddle())
	{
		edgexRouter.GET("/search", resp.JSONOutPutWrapper(edgex.SearchEdgex))
		edgexRouter.POST("/create", resp.JSONOutPutWrapper(edgex.CreateEdgex))
		edgexRouter.POST("/delete", resp.JSONOutPutWrapper(edgex.DeleteEdgex))
		edgexRouter.POST("/follow", resp.JSONOutPutWrapper(edgex.FollowEdgex))
		edgexRouter.POST("/unfollow", resp.JSONOutPutWrapper(edgex.UnFollowEdgex))
	}
	userRouter := r.Group("/edgex_admin/user")
	{
		userRouter.POST("/register", resp.JSONOutPutWrapper(user.Register))
		userRouter.POST("/login", session.SessionMiddleware(), resp.JSONOutPutWrapper(user.Login))
		userRouter.GET("/logout", resp.JSONOutPutWrapper(user.Logout))
	}
}
