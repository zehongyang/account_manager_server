package db

import (
	"github.com/zehongyang/bee/dbs"
	"github.com/zehongyang/bee/utils"
)

type DBUser struct {
	db *dbs.SplitTable
}

var GetDBUser = utils.Single(func() *DBUser {
	return &DBUser{db: dbs.Get("user")}
})
