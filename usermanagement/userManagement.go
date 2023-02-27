package usermanagement

import (
	"strings"
	"time"

	wapi "github.com/iamacarpet/go-win64api"
	"go.uber.org/zap"
	"golang.org/x/sys/windows/registry"
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

func ReturnUsers() []LocalUser {
	userslice := make([]LocalUser, 0)
	users1, err := wapi.ListLocalUsers()
	if err != nil {
		zap.S().Error("Error fetching user list, %s.\r\n", err.Error())

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

	//fmt.Printf("%v", userslice[0]) // Add api call to serial scripter api to compare original user list with current user list
	return userslice

}

func GetLastLoggenOnUser() string {
	registryKey := "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\LogonUI"

	var access uint32 = registry.QUERY_VALUE
	regKey, err := registry.OpenKey(registry.LOCAL_MACHINE, registryKey, access)
	if err != nil {
		if err != registry.ErrNotExist {
			zap.S().Error(err)
		}
		return ""
	}

	user, _, err := regKey.GetStringValue("LastLoggedOnUser")
	if err != nil {
		zap.S().Error(err)
		return ""
	}

	onlyUser := strings.TrimLeft(user, ".\\")

	return onlyUser
}

func noAdmin(user string) {
	_, err := wapi.RevokeAdmin(user)
	if err != nil {
		zap.S().Error(err)
	}
}

func disableUser(user string) {
	u := true
	_, err := wapi.UserDisabled(user, u)
	if err != nil {
		zap.S().Error(err)
	}

}

func enableUser(user string) {
	u := false
	_, err := wapi.UserDisabled(user, u)
	if err != nil {
		zap.S().Error(err)
	}

}

func delUser(user string) {
	_, err := wapi.UserDelete(user)
	if err != nil {
		zap.S().Error(err)
	}
}

func addUser(username, fullname, password string) {
	_, err := wapi.UserAdd(username, fullname, password)
	if err != nil {
		zap.S().Error(err)
	}
}

func usrNochangepw(user string) {
	c := true // password doesnt expire
	_, err := wapi.UserPasswordNoExpires(user, c)
	if err != nil {
		zap.S().Error(err)
	}
}

func usrChangepw(user string) {
	c := false //password expires
	_, err := wapi.UserPasswordNoExpires(user, c)
	if err != nil {
		zap.S().Error(err)
	}
}

func forcePasswdchange(user, newpasswd string) {
	_, err := wapi.ChangePassword(user, newpasswd)
	if err != nil {
		zap.S().Error(err)
	}
}

func setFullNameAttribute(user, fullname string) {
	_, err := wapi.UserUpdateFullname(user, fullname)
	if err != nil {
		zap.S().Error(err)
	}
}
