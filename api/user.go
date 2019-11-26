package main

import (
	"crypto/subtle"
	"time"
)

type userPrivilege struct {
	Download bool
	Upload   bool
	Delete   bool
}

type user struct {
	Expire    time.Time
	Password  string
	Privilege userPrivilege
}

type userList map[string]*user

func compare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func (u userList) Set(users confUsers) (err error) {
	for key := range u {
		delete(u, key)
	}
	for name, info := range users {
		var ex time.Time
		if info.Expire != "" {
			ex, err = time.ParseInLocation("2006-01-02 15:04:05", info.Expire, time.Local)
			if err != nil {
				return
			}
		}

		user := &user{Expire: ex, Password: info.Password}
		for _, pri := range info.Privilege {
			switch pri {
			case "download":
				user.Privilege.Download = true
			case "upload":
				user.Privilege.Upload = true
			case "delete":
				user.Privilege.Delete = true
			}
		}

		u[name] = user
	}
	return
}

func (u userList) verify(username, password string) (*user, bool) {
	var user *user
	var ok bool

	if username == "" || username == "anonymous" {
		if user, ok = u["anonymous"]; !ok {
			return nil, false
		}
	} else {
		if user, ok = u[username]; !ok || !compare(user.Password, password) {
			return nil, false
		}
	}

	if !user.Expire.IsZero() && time.Now().After(user.Expire) {
		return nil, false
	}

	return user, true
}
