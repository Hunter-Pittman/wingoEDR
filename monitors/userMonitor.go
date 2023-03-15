package monitors

import (
	"encoding/json"
	"os"
	"time"
	"wingoEDR/apis/serialscripter"
	"wingoEDR/db"
	"wingoEDR/usermanagement"

	"github.com/blockloop/scan/v2"
	"go.uber.org/zap"
)

// Payload fields:
// Last logged on user

type MonitorPayload struct {
	NewUserInfo      []usermanagement.User
	LastLoggedOnUser string
}

func InitUsers() {
	users := usermanagement.ReturnUsers()

	if db.CountTableRecords("currentusers") == 0 {
		usersToDB(users, false)
	}

}

func MonitorUsers() {
	users := usermanagement.ReturnUsers()
	queriedUsers := usersFromDB()

	if len(users) > len(queriedUsers) {
		for _, queriedUser := range queriedUsers {
			for i, user := range users {
				if user.Username == queriedUser.Username {
					users[i] = users[len(users)-1]
					users = users[:len(users)-1]
				}
			}
		}
		for _, user := range users {
			zap.S().Info("New user detected: ", user.Username)
			newUserIncident([]usermanagement.User{user})
		}

		usersToDB(users, false)
	} else if len(users) < len(queriedUsers) {
		for _, user := range users {
			for i, queriedUser := range queriedUsers {
				if user.Username == queriedUser.Username {
					queriedUsers[i] = queriedUsers[len(queriedUsers)-1]
					queriedUsers = queriedUsers[:len(queriedUsers)-1]
				}
			}
		}

		deleteUserFromDB(queriedUsers)
	} else {
		usersToDB(users, true)
	}
}

func usersToDB(users []usermanagement.User, update bool) {
	var operation string
	if update {
		operation = "UPDATE currentusers SET fullname = ?, enabled = ?, locked = ?, admin = ?, passwdexpired = ?, cantchangepasswd = ?, passwdage = ?, lastlogon = ?, badpasswdattempts = ?, numoflogons = ? WHERE username = ?"
	} else {
		operation = "INSERT INTO currentusers (username, fullname, enabled, locked, admin, passwdexpired, cantchangepasswd, passwdage, lastlogon, badpasswdattempts, numoflogons) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	}

	conn := db.DbConnect()
	stmt, err := conn.Prepare(operation)
	if err != nil {
		zap.S().Error("Error inserting user into database: ", err)
	}

	for _, user := range users {
		_, err := stmt.Exec(user.Username, user.Fullname, user.Enabled, user.Locked, user.Admin, user.PasswdExpired, user.CantChangePasswd, user.PasswdAge, user.LastLogon.String(), user.BadPasswdAttempts, user.NumOfLogons)
		if err != nil {
			zap.S().Error("Error inserting user into database: ", err)
		}
	}

	defer stmt.Close()
}

func usersFromDB() []usermanagement.User {
	conn := db.DbConnect()
	rows, err := conn.Query("SELECT * FROM currentusers")
	if err != nil {
		zap.S().Error("Error fetching user list from database: ", err)
	}

	var queriedUsers []usermanagement.User
	err1 := scan.Rows(&queriedUsers, rows)
	if err1 != nil {
		zap.S().Error("Error scanning user list from database: ", err1)
	}

	defer rows.Close()

	return queriedUsers
}

func deleteUserFromDB(user []usermanagement.User) {
	for _, user := range user {
		zap.S().Info("User deleted: ", user.Username)

		conn := db.DbConnect()
		stmt, err := conn.Prepare("DELETE FROM currentusers WHERE username = ?")
		if err != nil {
			zap.S().Error("Error deleting user from database: ", err)
		}

		_, err = stmt.Exec(user.Username)
		if err != nil {
			zap.S().Error("Error deleting user from database: ", err)
		}

		defer stmt.Close()
	}

}

func newUserIncident(newUsers []usermanagement.User) {

	payload := MonitorPayload{
		NewUserInfo:      newUsers,
		LastLoggedOnUser: usermanagement.GetLastLoggenOnUser(),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		zap.S().Error("Error marshalling payload: ", err)
	}

	newIncident := serialscripter.Incident{
		Name:        "New User Created",
		CurrentTime: time.Now().String(),
		User:        "New user detected",
		Severity:    "High",
		Payload:     string(jsonPayload),
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	newAlert := serialscripter.Alert{
		Host:     hostname,
		Incident: newIncident,
	}
	serialscripter.IncidentAlert(newAlert)
}
