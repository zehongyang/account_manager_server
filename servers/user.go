package servers

import (
	"account_manager/proto/pb"
	"fmt"
	"github.com/zehongyang/bee/config"
	"github.com/zehongyang/bee/logger"
	"github.com/zehongyang/bee/utils"
)

const (
	WxLoginUrl = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type AppWxConfig struct {
	WxConfig struct {
		AppId     string `json:"appid"`
		AppSecret string `json:"appsecret"`
	}
}

type UserLoginServer struct {
	cfg AppWxConfig
}

var GetUserLoginServer = utils.Single(func() *UserLoginServer {
	var wc AppWxConfig
	err := config.Load(&wc)
	if err != nil {
		logger.Fatal().Err(err).Msg("GetUserLoginServer")
	}
	return &UserLoginServer{cfg: wc}
})

func (s *UserLoginServer) Login(code string) (*pb.WxLoginPayload, error) {
	url := fmt.Sprintf(WxLoginUrl, s.cfg.WxConfig.AppId, s.cfg.WxConfig.AppSecret, code)
	var res pb.WxLoginPayload
	err := utils.Get(url, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
