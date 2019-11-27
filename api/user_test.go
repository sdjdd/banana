package main

import (
	"testing"
	"time"
)

func TestSetUser(t *testing.T) {
	users := userList{}
	users.set(confUsers{
		"user": {
			Expire:    "2019-01-01 20:00:00",
			Password:  "secret",
			Privilege: []string{"download", "upload", "delete"},
		},
	})

	u, ok := users["user"]
	switch {
	case !ok:
		t.Fatal("user should exists")
	case !u.Expire.Equal(time.Date(2019, 1, 1, 20, 0, 0, 0, time.Local)):
		t.Fatal("expire time parsing not correct")
	case u.Password != "secret":
		t.Fatal("password not correct")
	case u.Privilege.Download != true || u.Privilege.Upload != true ||
		u.Privilege.Delete != true:
		t.Fatal("privilege parsing not correct")
	}
}

func TestVerifyUser(t *testing.T) {
	users := userList{}
	users.set(confUsers{
		"user": {
			Password: "secret",
			Expire:   time.Now().Add(time.Hour).Format("2006-01-02 15:04:05"),
		},
		"anonymous": {},
	})

	if _, ok := users.verify("user", "secret"); !ok {
		t.Fatal("user should be passed")
	}
	if _, ok := users.verify("", "anything"); !ok {
		t.Fatal("anonymous should be passed")
	}
	if _, ok := users.verify("stranger", "anything"); ok {
		t.Fatal("stranger should be rejected")
	}

	users["anonymous"].Expire = time.Now().Add(-time.Hour)
	if _, ok := users.verify("anonymous", "anything"); ok {
		t.Fatal("expired anonymous should be rejected")
	}

	delete(users, "anonymous")
	if _, ok := users.verify("anonymous", "anything"); ok {
		t.Fatal("unset anonymous should be rejected")
	}
}
