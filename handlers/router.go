package handlers

import (
	"github.com/zehongyang/bee"
)

func InitRouter(srv *bee.HttpServer) {
	if srv == nil {
		return
	}
	srv.Post("/login", UserLoginQuery())
	srv.Post("/upload", FileUploadQuery())
}
