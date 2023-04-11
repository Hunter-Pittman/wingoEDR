package usermanagement

import (
	"strings"
	"time"

	wapi "github.com/iamacarpet/go-win64api"
	"go.uber.org/zap"
	"golang.org/x/sys/windows/registry"
)

type User struct {
	Username          string
	Fullname          string
	Enabled           bool
	Locked            bool
	Admin             bool
	PasswdExpired     bool
	CantChangePasswd  bool
	PasswdAge         time.Duration
	LastLogon         time.Time
	BadPasswdAttempts uint32
	NumOfLogons       uint32
}

func ReturnUsers() []User {
	userslice := make([]User, 0)
	users1, err := wapi.ListLocalUsers()
	if err != nil {
		zap.S().Error("Error fetching user list, %s.\r\n", err.Error())

	}
	for _, u := range users1 {
		userlist := User{
			Username:          u.Username,
			Fullname:          u.FullName,
			Enabled:           u.IsEnabled,
			Locked:            u.IsLocked,
			Admin:             u.IsAdmin,
			PasswdExpired:     u.PasswordNeverExpires,
			CantChangePasswd:  u.NoChangePassword,
			PasswdAge:         u.PasswordAge,
			LastLogon:         u.LastLogon,
			BadPasswdAttempts: u.BadPasswordCount,
			NumOfLogons:       u.NumberOfLogons}

		userslice = append(userslice, userlist)

	}

	//fmt.Printf("Userslice: %v", userslice)
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

func EnableUser(user string) {
	u := false
	_, err := wapi.UserDisabled(user, u)
	if err != nil {
		zap.S().Error(err)
	}

}

func DelUser(user string) {
	_, err := wapi.UserDelete(user)
	if err != nil {
		zap.S().Error(err)
	}
}

func AddUser(username, fullname, password string) error {
	_, err := wapi.UserAdd(username, fullname, password)
	if err != nil {
		zap.S().Error("Error adding user: ", err)
		return err
	}

	return nil
}

func UsrNochangepw(user string) {
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

func ForcePasswdchange(user, newpasswd string) {
	_, err := wapi.ChangePassword(user, newpasswd)
	if err != nil {
		zap.S().Error(err)
	}
}

func SetFullNameAttribute(user, fullname string) {
	_, err := wapi.UserUpdateFullname(user, fullname)
	if err != nil {
		zap.S().Error(err)
	}
}
