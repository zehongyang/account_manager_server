package handlers

import "github.com/zehongyang/bee"

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
