package main

import (
	"fmt"
	"time"

	wapi "github.com/iamacarpet/go-win64api"
)

type LocalUser struct {
	Username          string
	Fullname          string
	Enabled           bool
	Locked            bool
	Admin             bool
	Passwdexpired     bool
	CantChangePasswd  bool
	Passwdage         time.Duration
	Lastlogon         time.Time
	BadPasswdAttempts uint32
	NumofLogons       uint32
}

func main() {
	returnUsers()
}

func returnUsers() []LocalUser {
	userslice := make([]LocalUser, 0)
	users1, err := wapi.ListLocalUsers()
	if err != nil {
		fmt.Printf("Error fetching user list, %s.\r\n", err.Error())

	}
	for _, u := range users1 {
		userlist := LocalUser{
			Username:          u.Username,
			Fullname:          u.FullName,
			Enabled:           u.IsEnabled,
			Locked:            u.IsLocked,
			Admin:             u.IsAdmin,
			Passwdexpired:     u.PasswordNeverExpires,
			CantChangePasswd:  u.NoChangePassword,
			Passwdage:         u.PasswordAge,
			Lastlogon:         u.LastLogon,
			BadPasswdAttempts: u.BadPasswordCount,
			NumofLogons:       u.NumberOfLogons}

		userslice = append(userslice, userlist)

	}

	fmt.Printf("%v", userslice[0])
	return userslice

}
