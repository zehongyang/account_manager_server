package handlers

import (
	"account_manager/db"
	"account_manager/models"
	"account_manager/proto/pb"
	"account_manager/servers"
	"github.com/zehongyang/bee"
	"github.com/zehongyang/bee/logger"
	"github.com/zehongyang/bee/utils"
	"net/http"
)

func UserLoginQuery() bee.Handler {
	srv := servers.GetUserLoginServer()
	dbUser := db.GetDBUser()
	return func(req bee.IContext) {
		var (
			q   pb.UserLoginQuery
			res pb.UserLoginQueryResponse
		)
		err := req.Bind(&q)
		if err != nil || len(q.Code) < 1 {
			logger.Error().Err(err).Any("q", &q).Msg("UserLoginQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		rs, err := srv.Login(q.Code)
		if err != nil {
			logger.Error().Err(err).Any("q", &q).Msg("UserLoginQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		if len(rs.Openid) < 1 {
			logger.Error().Err(err).Any("q", &q).Msg("UserLoginQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		var token = srv.GenerateToken()
		var tm = utils.Now()
		var ip = req.GetIp()
		var user = &models.User{
			Openid: rs.Openid,
			Token:  token,
			Ctm:    tm,
			Ltm:    tm,
			Cip:    ip,
			Lip:    ip,
		}
		err = dbUser.Upsert(user)
		if err != nil {
			logger.Error().Err(err).Any("q", &q).Msg("UserLoginQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		res.Token = token
		res.User = &pb.UserInfo{
			Id:       int64(user.Id),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		}
		req.ResponseOk(&res)
	}
}
