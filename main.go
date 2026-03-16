package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zehongyang/bee"
	"github.com/zehongyang/bee/logger"
)

func main() {
	srv := bee.NewHttpServer()
	srv.Post("/test", func(ctx bee.IContext) {
		ctx.ResponseOk(gin.H{"name": "zhangsan"})
	})
	err := srv.Run(":22345")
	if err != nil {
		logger.Error().Err(err).Msg("running http server")
	}
}
