package handlers

import (
	"github.com/zehongyang/bee"
)

func InitRouter(srv *bee.HttpServer) {
	if srv == nil {
		return
	}
	srv.Post("/ac/login", UserLoginQuery())
	authGroup := srv.Group("/ac")
	authGroup.Use(Auth())
	authGroup.Post("/upload", FileUploadQuery())
	authGroup.Post("/user/modify", UserInfoModifyQuery())
	authGroup.Post("/key/set", MasterKeySetQuery())
	authGroup.Post("/key/verify", MasterKeyVerifyQuery())
	authGroup.Post("/app/add", AppAddQuery())
	authGroup.Post("/app/list", AppListQuery())
}
