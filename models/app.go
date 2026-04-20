package models

type App struct {
	Id      int    `xorm:"not null pk autoincr INTEGER"`
	Uid     int    `xorm:"not null index(app_uik) INTEGER"`
	Name    string `xorm:"index(app_uik) VARCHAR(255)"`
	Payload string `xorm:"TEXT"`
	Ctm     int64  `xorm:"BIGINT"`
}

func (m *App) TableName() string {
	return "app"
}
