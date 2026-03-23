package models

type User struct {
	Openid   string `xorm:"not null pk comment('wx open id') VARCHAR(40)"`
	Unionid  string `xorm:"comment('wx union id') VARCHAR(40)"`
	Nickname string `xorm:"comment('昵称') VARCHAR(40)"`
	Avatar   string `xorm:"comment('头像') VARCHAR(255)"`
	Token    string `xorm:"comment('登录token') VARCHAR(255)"`
	Ctm      int64  `xorm:"comment('创建时间') BIGINT"`
	Ltm      int64  `xorm:"comment('最近登录时间') BIGINT"`
	Lip      string `xorm:"comment('最近登录ip') VARCHAR(255)"`
	Cip      string `xorm:"comment('注册ip') VARCHAR(255)"`
}

func (m *User) TableName() string {
	return "user"
}
