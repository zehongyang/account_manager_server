package db

import (
	"account_manager/models"
	"time"

	"github.com/zehongyang/bee/caches"
	"github.com/zehongyang/bee/dbs"
	"github.com/zehongyang/bee/utils"
)

type DBUser struct {
	db    *dbs.SplitTable
	cache *caches.Cache[int64, *models.User]
}

var GetDBUser = utils.Single(func() *DBUser {
	return &DBUser{db: dbs.Get("user"), cache: caches.NewCache[int64, *models.User](1000, time.Hour)}
})

func (s *DBUser) Upsert(v *models.User) error {
	sql := "INSERT INTO `user` (openid,token,ctm,ltm,lip,cip) VALUES(?,?,?,?,?,?) ON CONFLICT(openid) DO UPDATE SET ltm=?,lip=? RETURNING *"
	_, err := s.db.Num(0).SQL(sql, v.Openid, v.Token, v.Ctm, v.Ltm, v.Lip, v.Cip, v.Ltm, v.Lip).Get(v)
	return err
}

func (s *DBUser) GetUser(uid int64) (*models.User, error) {
	var u models.User
	has, err := s.db.Num(uid).Where("id=?", uid).Get(&u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &u, nil
}

func (s *DBUser) GetUserCache(uid int64) (*models.User, error) {
	value, ok := s.cache.Get(uid)
	if ok {
		return value, nil
	}
	user, err := s.GetUser(uid)
	if err != nil {
		return nil, err
	}
	s.cache.Add(uid, user)
	return user, nil
}

func (s *DBUser) Modify(v models.User, cols ...string) (int64, error) {
	ses := s.db.Num(int64(v.Id)).Where("id=?", v.Id)
	if len(cols) > 0 {
		ses.Cols(cols...)
	}
	return ses.Update(v)
}
