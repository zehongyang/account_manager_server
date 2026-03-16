package servers

import (
	"account_manager/proto/pb"
	"fmt"
	"github.com/zehongyang/bee/utils"
)

const (
	WxLoginUrl = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	AppId      = ""
	AppSecret  = ""
)

type UserLoginServer struct {
}

var GetUserLoginServer = utils.Single(func() *UserLoginServer {
	return &UserLoginServer{}
})

func (s *UserLoginServer) Login(code string) (*pb.WxLoginPayload, error) {
	url := fmt.Sprintf(WxLoginUrl, AppId, AppSecret, code)
	var res pb.WxLoginPayload
	err := utils.Get(url, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
