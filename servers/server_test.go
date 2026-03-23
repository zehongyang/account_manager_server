package servers

import "testing"

func TestWxCodeLogin(t *testing.T) {
	srv := GetUserLoginServer()
	res, err := srv.Login("0f1YyqFa1S9VmL0qmBHa15JC6d1YyqFp")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestEnv(t *testing.T) {
	srv := GetUserLoginServer()
	t.Log(srv.cfg)
}
