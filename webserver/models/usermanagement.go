package models

import (
	"wingoEDR/usermanagement"
)

type NewUser struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
}

func CreateNewUser(user NewUser) error {
	err := usermanagement.AddUser(user.Username, user.Fullname, user.Password)
	if err != nil {
		return err
	}
	return nil
}
