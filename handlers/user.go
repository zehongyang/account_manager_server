package handlers

import (
	"account_manager/db"
	"account_manager/models"
	"account_manager/proto/pb"
	"account_manager/servers"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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
	dbApp := db.GetDBApp()
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
		var amp = make(map[int64]string)
		if len(q.Apps) > 0 {
			for _, app := range q.Apps {
				if app.Id < 1 || app.Payload == nil || len(app.Payload.AccountPayload) < 1 ||
					len(app.Payload.PwdPayload) < 1 {
					logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("MasterKeySetQuery")
					req.ResponseError(http.StatusBadRequest)
					return
				}
				marshal, err := json.Marshal(app.Payload)
				if err != nil {
					logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("MasterKeySetQuery")
					req.ResponseError(http.StatusInternalServerError)
					return
				}
				amp[app.Id] = string(marshal)
			}
		}
		_, err = dbUser.Modify(models.User{Payload: q.Payload, Ciphertext: q.Ciphertext, Id: int(uid)}, "payload", "ciphertext")
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("MasterKeySetQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		if len(q.Apps) > 0 {
			for _, app := range q.Apps {
				_, err = dbApp.Modify(models.App{
					Id:      int(app.Id),
					Uid:     int(uid),
					Payload: amp[app.Id],
				}, "payload")
				if err != nil {
					logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("MasterKeySetQuery")
					req.ResponseError(http.StatusInternalServerError)
					return
				}
			}
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

func AppAddQuery() bee.Handler {
	dbApp := db.GetDBApp()
	return func(req bee.IContext) {
		var (
			q   pb.AppAddQuery
			res pb.AppAddQueryResponse
		)
		uid := req.GetAccount().GetUid()
		err := req.Bind(&q)
		if err != nil || len(q.Name) < 1 || q.Payload == nil || len(q.Payload.AccountPayload) < 1 ||
			len(q.Payload.PwdPayload) < 1 {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppAddQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		count, err := dbApp.Count(uid)
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppAddQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		if count >= 30 {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppAddQuery")
			req.ResponseError(http.StatusBadRequest, "应用数量已达上限")
			return
		}
		exists, err := dbApp.Exists(uid, q.Name)
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppAddQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		if exists {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppAddQuery")
			req.ResponseError(http.StatusBadRequest, "应用已存在")
			return
		}
		strPayload, err := json.Marshal(q.Payload)
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppAddQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		_, err = dbApp.Create(models.App{
			Name:    q.Name,
			Uid:     int(uid),
			Payload: string(strPayload),
			Ctm:     utils.Now(),
		})
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppAddQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		req.ResponseOk(&res)
	}
}

func AppListQuery() bee.Handler {
	dbApp := db.GetDBApp()
	return func(req bee.IContext) {
		var (
			q   pb.AppListQuery
			res pb.AppListQueryResponse
		)
		uid := req.GetAccount().GetUid()
		err := req.Bind(&q)
		if err != nil || q.Id < 1 {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppListQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		apps, err := dbApp.List(uid, q.Id, 5, q.All)
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppListQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		for _, app := range apps {
			var pai = &pb.AppInfo{
				Id:   int64(app.Id),
				Name: app.Name,
			}
			if len(app.Payload) > 0 {
				var payload pb.AppPayload
				err = json.Unmarshal([]byte(app.Payload), &payload)
				if err != nil {
					logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppListQuery")
					req.ResponseError(http.StatusInternalServerError)
					return
				}
				pai.Payload = &payload
			}
			res.Lists = append(res.Lists, pai)
		}
		req.ResponseOk(&res)
	}
}

func AppModifyQuery() bee.Handler {
	dbApp := db.GetDBApp()
	return func(req bee.IContext) {
		var (
			q   pb.AppModifyQuery
			res pb.AppModifyQueryResponse
		)
		uid := req.GetAccount().GetUid()
		err := req.Bind(&q)
		if err != nil || q.Id < 1 || len(q.Name) < 1 || q.Payload == nil ||
			len(q.Payload.AccountPayload) < 1 || len(q.Payload.PwdPayload) < 1 {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppModifyQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		marshal, err := json.Marshal(q.Payload)
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppModifyQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		_, err = dbApp.Modify(models.App{
			Id:      int(q.Id),
			Uid:     int(uid),
			Name:    q.Name,
			Payload: string(marshal),
		})
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppModifyQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		req.ResponseOk(&res)
	}
}

func AppRemoveQuery() bee.Handler {
	dbApp := db.GetDBApp()
	return func(req bee.IContext) {
		var (
			q   pb.AppRemoveQuery
			res pb.AppRemoveQueryResponse
		)
		uid := req.GetAccount().GetUid()
		err := req.Bind(&q)
		if err != nil || q.Id < 1 {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppRemoveQuery")
			req.ResponseError(http.StatusBadRequest)
			return
		}
		_, err = dbApp.Remove(q.Id, uid)
		if err != nil {
			logger.Error().Err(err).Any("uid", uid).Any("q", &q).Msg("AppRemoveQuery")
			req.ResponseError(http.StatusInternalServerError)
			return
		}
		req.ResponseOk(&res)
	}
}
