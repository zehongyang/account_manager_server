package handlers

import (
	"account_manager/proto/pb"
	"account_manager/servers"
	"github.com/zehongyang/bee"
	"github.com/zehongyang/bee/logger"
	"net/http"
)

func UserLoginQuery() bee.Handler {
	srv := servers.GetUserLoginServer()
	return func(req bee.IContext) {
		var (
			q   pb.UserLoginQuery
			res pb.UserLoginQueryResponse
		)
		err := req.Bind(&q)
		if err != nil || len(q.Code) < 1 {
			logger.Error().Err(err).Msg("UserLoginQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		srv.Login(q.Code)
		req.ResponseOk(&res)
	}
}
