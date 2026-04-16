package handlers

import (
	"account_manager/db"
	"account_manager/models"
	"account_manager/proto/pb"
	"account_manager/servers"
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"github.com/zehongyang/bee"
	"github.com/zehongyang/bee/logger"
	"github.com/zehongyang/bee/utils"
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
		res.Token = user.Token
		var vefifyData string
		if len(user.Ciphertext) > 0 {
			sum := md5.Sum([]byte(user.Ciphertext))
			vefifyData = hex.EncodeToString(sum[:])
		}
		res.User = &pb.UserInfo{
			Id:       int64(user.Id),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Payload:  user.Payload,
			Key:      vefifyData,
		}
		req.ResponseOk(&res)
	}
}

func UserInfoModifyQuery() bee.Handler {
	dbUser := db.GetDBUser()
	return func(req bee.IContext) {
		var (
			q   pb.UserInfoModifyQuery
			res pb.UserInfoModifyQueryResponse
		)
		uid := req.GetAccount().GetUid()
		err := req.Bind(&q)
		if err != nil || len(q.Avatar) < 1 || len(q.Nickname) < 1 {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("UserInfoModifyQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		_, err = dbUser.Modify(models.User{
			Nickname: q.Nickname,
			Avatar:   q.Avatar,
			Id:       int(uid),
		})
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("UserInfoModifyQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		req.ResponseOk(&res)
	}
}

func MasterKeySetQuery() bee.Handler {
	dbUser := db.GetDBUser()
	return func(req bee.IContext) {
		var (
			q   pb.MasterKeySetQuery
			res pb.MasterKeySetQueryResponse
		)
		uid := req.GetAccount().GetUid()
		err := req.Bind(&q)
		if err != nil || len(q.Payload) < 1 || len(q.Ciphertext) < 1 {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("MasterKeySetQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		_, err = dbUser.Modify(models.User{Payload: q.Payload, Ciphertext: q.Ciphertext, Id: int(uid)}, "payload", "ciphertext")
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("MasterKeySetQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		sum := md5.Sum([]byte(q.Ciphertext))
		res.Key = hex.EncodeToString(sum[:])
		req.ResponseOk(&res)
	}
}

func MasterKeyVerifyQuery() bee.Handler {
	dbUser := db.GetDBUser()
	return func(req bee.IContext) {
		var (
			q   pb.MasterKeyVerifyQuery
			res pb.MasterKeyVerifyQueryResponse
		)
		uid := req.GetAccount().GetUid()
		err := req.Bind(&q)
		if err != nil || len(q.Ciphertext) < 1 {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("MasterKeyVerifyQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		user, err := dbUser.GetUser(uid)
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("MasterKeyVerifyQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		if user == nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("MasterKeyVerifyQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		res.Ok = q.Ciphertext == user.Ciphertext
		if res.Ok {
			sum := md5.Sum([]byte(q.Ciphertext))
			res.Key = hex.EncodeToString(sum[:])
		}
		req.ResponseOk(&res)
	}
}
