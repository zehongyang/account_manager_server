package db

import (
	"account_manager/models"
	"testing"
)

func TestCreateUser(t *testing.T) {
	dbUser := GetDBUser()
	var u = models.User{Openid: "1111"}
	err := dbUser.Upsert(&u)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}
