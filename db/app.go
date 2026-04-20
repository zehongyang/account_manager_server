package db

import (
	"account_manager/models"

	"github.com/zehongyang/bee/dbs"
	"github.com/zehongyang/bee/utils"
)

type DBApp struct {
	db *dbs.SplitTable
}

var GetDBApp = utils.Single(func() *DBApp {
	return &DBApp{db: dbs.Get("app")}
})

func (s *DBApp) Create(v models.App) (int64, error) {
	return s.db.Num(0).Insert(v)
}

func (s *DBApp) Exists(uid int64, name string) (bool, error) {
	return s.db.Num(0).Where("uid = ? and name = ?", uid, name).Exist()
}

func (s *DBApp) List(uid, id int64, size int) ([]models.App, error) {
	var apps []models.App
	err := s.db.Num(0).Where("uid = ? and id < ?", uid, id).Limit(size).Desc("id").Find(&apps)
	return apps, err
}
