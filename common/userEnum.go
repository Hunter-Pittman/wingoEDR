package common

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

func common() {
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

func noAdmin(user string) {
	_, err := wapi.RevokeAdmin(user)
	if err != nil {
		fmt.Println(err)
	}
}

func disableUser(user string) {
	u := true
	_, err := wapi.UserDisabled(user, u)
	if err != nil {
		fmt.Println(err)
	}

}

func enableUser(user string) {
	u := false
	_, err := wapi.UserDisabled(user, u)
	if err != nil {
		fmt.Println(err)
	}

}

func delUser(user string) {
	_, err := wapi.UserDelete(user)
	if err != nil {
		fmt.Println(err)
	}
}

func addUser(username, fullname, password string) {
	_, err := wapi.UserAdd(username, fullname, password)
	if err != nil {
		fmt.Println(err)
	}
}

func usrNochangepw(user string) {
	c := true // password doesnt expire
	_, err := wapi.UserPasswordNoExpires(user, c)
	if err != nil {
		fmt.Println(err)
	}
}

func usrChangepw(user string) {
	c := false //password expires
	_, err := wapi.UserPasswordNoExpires(user, c)
	if err != nil {
		fmt.Println(err)
	}
}

func forcePasswdchange(user, newpasswd string) {
	_, err := wapi.ChangePassword(user, newpasswd)
	if err != nil {
		fmt.Println(err)
	}
}

func setFullNameAttribute(user, fullname string) {
	_, err := wapi.UserUpdateFullname(user, fullname)
	if err != nil {
		fmt.Println(err)
	}
}
