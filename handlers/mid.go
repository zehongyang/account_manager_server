package handlers

import (
	"account_manager/db"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/zehongyang/bee"
	"github.com/zehongyang/bee/logger"
	"github.com/zehongyang/bee/utils"
	"net/http"
)

type Header struct {
	Uid       int64  `header:"Uid"`
	Token     string `header:"Token"`
	Timestamp int64  `header:"Timestamp"`
	FileHash  string `header:"FileHash"`
}

func Cors() bee.Handler {
	return func(c bee.IContext) {
		method := c.GetMethod()

		c.SetHeader("Access-Control-Allow-Origin", "*")
		c.SetHeader("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.SetHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.SetHeader("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.SetHeader("Access-Control-Allow-Credentials", "true")

		// 放行所有 OPTIONS 方法
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Auth() bee.Handler {
	dbUser := db.GetDBUser()
	return func(c bee.IContext) {
		var h Header
		err := c.BindHeader(&h)
		if err != nil || h.Uid < 1 || h.Token == "" {
			logger.Error().Err(err).Any("h", h).Msg("auth")
			c.ResponseError(http.StatusBadRequest)
			return
		}
		if utils.Abs(utils.Now()-h.Timestamp) > utils.Time1MinMs {
			logger.Error().Err(err).Any("h", h).Msg("auth")
			c.ResponseError(http.StatusRequestTimeout)
			return
		}
		user, err := dbUser.GetUser(h.Uid)
		if err != nil {
			logger.Error().Err(err).Any("h", h).Msg("auth")
			c.ResponseError(http.StatusInternalServerError)
			return
		}
		if user == nil {
			logger.Error().Err(err).Any("h", h).Msg("auth")
			c.ResponseError(http.StatusBadRequest)
			return
		}
		sum := md5.Sum([]byte(fmt.Sprintf("%s%d%d", user.Token, h.Uid, h.Timestamp)))
		if hex.EncodeToString(sum[:]) != h.Token {
			logger.Error().Err(err).Any("h", h).Msg("auth")
			c.ResponseError(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
