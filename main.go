package main

import (
	"account_manager/handlers"
	"github.com/zehongyang/bee"
	"github.com/zehongyang/bee/logger"
)

func main() {
	srv := bee.NewHttpServer()
	handlers.InitRouter(srv)
	err := srv.Run(":22345")
	if err != nil {
		logger.Error().Err(err).Msg("running http server")
	}
}
