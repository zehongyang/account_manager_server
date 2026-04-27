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

func (s *DBApp) Count(uid int64) (int64, error) {
	return s.db.Num(0).Where("uid=?", uid).Count()
}

func (s *DBApp) Create(v models.App) (int64, error) {
	return s.db.Num(0).Insert(v)
}

func (s *DBApp) Exists(uid int64, name string) (bool, error) {
	return s.db.Num(0).Where("uid = ? and name = ?", uid, name).Exist()
}

func (s *DBApp) List(uid, id int64, size int, all bool) ([]models.App, error) {
	var apps []models.App
	ses := s.db.Num(0).Where("uid = ? and id < ?", uid, id)
	if !all {
		ses.Limit(size)
	}
	err := ses.Desc("id").Find(&apps)
	return apps, err
}

func (s *DBApp) Modify(v models.App, cols ...string) (int64, error) {
	ses := s.db.Num(int64(v.Id)).Where("id=? and uid=?", v.Id, v.Uid)
	if len(cols) > 0 {
		ses.Cols(cols...)
	}
	return ses.Update(v)
}

func (s *DBApp) Remove(id, uid int64) (int64, error) {
	return s.db.Num(id).Where("id=? and uid=?", id, uid).Delete()
}
