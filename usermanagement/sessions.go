package usermanagement

import (
	"fmt"
	"wingoEDR/config"

	wapi "github.com/iamacarpet/go-win64api"
	"github.com/iamacarpet/go-win64api/shared"
	"go.uber.org/zap"
)

type SessionDetailsTermnationStatus struct {
	shared.SessionDetails
	SessionTerminated bool
}

func ListSessions() []shared.SessionDetails {
	// This check runs best as NT AUTHORITY\SYSTEM
	//
	// Running as a normal or even elevated user,
	// we can't properly detect who is an admin or not.
	//
	// This is because we require TOKEN_DUPLICATE permission,
	// which we don't seem to have otherwise (Win10).
	sessions, err := wapi.ListLoggedInUsers()
	if err != nil {
		zap.S().Error("Error getting session details: ", err)
	}

	whitelist := config.GetWhitelistedUsers()

	for _, u := range whitelist {
		for _, s := range sessions {
			if s.Username == u {
				fmt.Println("This was triggered")
				//sessions = RemoveIndex(sessions, i) // We still wanna include the whitelisted sessions, but need to
			}
		}
	}

	return sessions
}

func RemoveIndex(slice []shared.SessionDetails, index int) []shared.SessionDetails {
	return append(slice[:index], slice[index+1:]...)
}
