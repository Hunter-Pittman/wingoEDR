package monitors

import (
	"encoding/json"
	"os"
	"time"
	"wingoEDR/apis/serialscripter"
	"wingoEDR/db"
	"wingoEDR/shares"
	"wingoEDR/usermanagement"

	"github.com/blockloop/scan/v2"
	"go.uber.org/zap"
)

type ShareMonitorPayload struct {
	NewShareInfo     []shares.SMBInfo
	LastLoggedOnUser string
}

func InitShares() {
	shares := shares.GetShares()

	if db.CountTableRecords("currentshares") == 0 {
		sharesToDB(shares, false)
	}
}

func SharesMonitor() {
	localShares := shares.GetShares()
	queriedShares := sharesFromDB()

	if len(localShares) > len(queriedShares) {
		for _, queriedShare := range queriedShares {
			for i, share := range localShares {
				if share.NetName == queriedShare.NetName {
					localShares[i] = localShares[len(localShares)-1]
					localShares = localShares[:len(localShares)-1]
				}
			}
		}
		for _, share := range localShares {
			zap.S().Info("New share detected: ", share.NetName)
			newShareIncident([]shares.SMBInfo{share})
		}

		sharesToDB(localShares, false)
	} else if len(localShares) < len(queriedShares) {
		for _, share := range localShares {
			for i, queriedShare := range queriedShares {
				if share.NetName == queriedShare.NetName {
					queriedShares[i] = queriedShares[len(queriedShares)-1]
					queriedShares = queriedShares[:len(queriedShares)-1]
				}
			}
		}

		deleteShareFromDB(queriedShares)
	} else {
		sharesToDB(localShares, true)
	}
}

func sharesToDB(shares []shares.SMBInfo, update bool) {
	var operation string
	if update {
		operation = "UPDATE currentshares SET netname = ?, remark = ?, path = ?, type = ?, permissions = ?, maxuses = ?, currentuses = ?"
	} else {
		operation = "INSERT INTO currentshares (netname, remark, path, type, permissions, maxuses, currentuses) VALUES (?, ?, ?, ?, ?, ?, ?)"
	}

	conn := db.DbConnect()
	stmt, err := conn.Prepare(operation)
	if err != nil {
		zap.S().Error("Error inserting user into database: ", err)
	}

	for _, share := range shares {
		_, err := stmt.Exec(share.NetName, share.Remark, share.Path, share.Type, share.Permissions, share.MaxUses, share.CurrentUses)
		if err != nil {
			zap.S().Error("Error inserting share into database: ", err)
		}
	}

	defer stmt.Close()
}

func sharesFromDB() []shares.SMBInfo {
	conn := db.DbConnect()
	rows, err := conn.Query("SELECT * FROM currentshares")
	if err != nil {
		zap.S().Error("Error fetching share list from database: ", err)
	}

	var queriedShares []shares.SMBInfo
	err1 := scan.Rows(&queriedShares, rows)
	if err1 != nil {
		zap.S().Error("Error scanning share list from database: ", err1)
	}

	defer rows.Close()

	return queriedShares
}

func deleteShareFromDB(shares []shares.SMBInfo) {
	for _, share := range shares {
		zap.S().Info("Share deleted: ", share.NetName)

		conn := db.DbConnect()
		stmt, err := conn.Prepare("DELETE FROM currentshares WHERE netname = ?")
		if err != nil {
			zap.S().Error("Error deleting share from database: ", err)
		}

		_, err = stmt.Exec(share.NetName)
		if err != nil {
			zap.S().Error("Error deleting share from database: ", err)
		}

		defer stmt.Close()
	}

}

func newShareIncident(newShare []shares.SMBInfo) {

	payload := ShareMonitorPayload{
		NewShareInfo:     newShare,
		LastLoggedOnUser: usermanagement.GetLastLoggenOnUser(),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		zap.S().Error("Error marshalling payload: ", err)
	}

	newIncident := serialscripter.Incident{
		Name:        "New share Created",
		CurrentTime: time.Now().String(),
		User:        "New share detected",
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
