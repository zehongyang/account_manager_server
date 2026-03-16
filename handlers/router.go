package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zehongyang/bee"
)

func InitRouter(srv *bee.HttpServer) {
	if srv == nil {
		return
	}
	srv.Post("/test", func(ctx bee.IContext) {
		ctx.ResponseOk(gin.H{"name": "zhangsan"})
	})
}
