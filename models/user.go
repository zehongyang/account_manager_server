package models

type User struct {
	Openid     string `xorm:"not null comment('wx open id') unique VARCHAR(40)"`
	Unionid    string `xorm:"comment('wx union id') VARCHAR(40)"`
	Nickname   string `xorm:"comment('昵称') VARCHAR(40)"`
	Avatar     string `xorm:"comment('头像') VARCHAR(255)"`
	Token      string `xorm:"comment('登录token') VARCHAR(255)"`
	Ctm        int64  `xorm:"comment('创建时间') BIGINT"`
	Ltm        int64  `xorm:"comment('最近登录时间') BIGINT"`
	Lip        string `xorm:"comment('最近登录ip') VARCHAR(255)"`
	Cip        string `xorm:"comment('注册ip') VARCHAR(255)"`
	Id         int    `xorm:"not null pk autoincr INTEGER"`
	Payload    string `xorm:"comment('验证密码的加密数据') TEXT"`
	Ciphertext string `xorm:"comment('待对比的明文数据') VARCHAR(255)"`
}

func (m *User) TableName() string {
	return "user"
}
