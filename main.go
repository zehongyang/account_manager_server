package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zehongyang/bee"
)

func main() {
	srv := bee.NewHttpServer()
	srv.Post("/test", func(ctx bee.IContext) {
		ctx.ResponseOk(gin.H{"name": "zhangsan"})
	})
	err := srv.Run(":22345")
	fmt.Println(err)
}
