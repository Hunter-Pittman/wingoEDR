package monitors

import (
	"encoding/json"
	"os"
	"time"
	"wingoEDR/apis/serialscripter"
	"wingoEDR/autoruns"
	"wingoEDR/db"
	"wingoEDR/usermanagement"

	"github.com/blockloop/scan/v2"
	"go.uber.org/zap"
)

// Payload fields:
// Last logged on user

type AutorunsMonitorPayload struct {
	NewAutorunsInfo  []autoruns.AutorunsInfo
	LastLoggedOnUser string
}

func InitAutoruns() {
	autoruns := autoruns.GetAutoruns()

	if db.CountTableRecords("currentautoruns") == 0 {
		autorunsToDB(autoruns, false)
	}

}

func MonitorAutoruns() {
	monitoredAutoruns := autoruns.GetAutoruns()
	queriedAutoruns := autorunsFromDB()

	if len(monitoredAutoruns) > len(queriedAutoruns) {
		for _, queriedAutoruns := range queriedAutoruns {
			for i, autorun := range monitoredAutoruns {
				if autorun.Type == queriedAutoruns.Type {
					monitoredAutoruns[i] = monitoredAutoruns[len(monitoredAutoruns)-1]
					monitoredAutoruns = monitoredAutoruns[:len(monitoredAutoruns)-1]
				}
			}
		}
		for _, autorun := range monitoredAutoruns {
			zap.S().Info("New autorun detected: ", autorun.Type)
			newAutorunIncident([]autoruns.AutorunsInfo{autorun})
		}

		autorunsToDB(monitoredAutoruns, false)
	} else if len(monitoredAutoruns) < len(queriedAutoruns) {
		for _, autorun := range monitoredAutoruns {
			for i, queriedAutorun := range queriedAutoruns {
				if autorun.Type == queriedAutorun.Type {
					queriedAutoruns[i] = queriedAutoruns[len(queriedAutoruns)-1]
					queriedAutoruns = queriedAutoruns[:len(queriedAutoruns)-1]
				}
			}
		}

		deleteAutorunFromDB(queriedAutoruns)
	} else {
		autorunsToDB(monitoredAutoruns, true)
	}
}

func autorunsToDB(autoruns []autoruns.AutorunsInfo, update bool) {
	var operation string
	if update {
		operation = "UPDATE currentautoruns SET type = ?, location = ?, image_path = ?, image_name = ?, arguments = ?, md5 = ?, sha1 = ?, sha256 = ?"
	} else {
		operation = "INSERT INTO currentautoruns (type, location, image_path, image_name, arguments, md5, sha1, sha256) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	}

	conn := db.DbConnect()
	stmt, err := conn.Prepare(operation)
	if err != nil {
		zap.S().Error("Error inserting autorun into database: ", err)
	}

	for _, autoruns := range autoruns {
		_, err := stmt.Exec(autoruns.Type, autoruns.Location, autoruns.ImagePath, autoruns.ImageName, autoruns.Arguments, autoruns.MD5, autoruns.SHA1, autoruns.SHA256)
		if err != nil {
			zap.S().Error("Error inserting autorun into database: ", err)
		}
	}

	defer stmt.Close()
}

func autorunsFromDB() []autoruns.AutorunsInfo {
	conn := db.DbConnect()
	rows, err := conn.Query("SELECT * FROM currentautoruns")
	if err != nil {
		zap.S().Error("Error fetching autoruns list from database: ", err)
	}

	var queriedAutoruns []autoruns.AutorunsInfo
	err1 := scan.Rows(&queriedAutoruns, rows)
	if err1 != nil {
		zap.S().Error("Error scanning autoruns list from database: ", err1)
	}

	defer rows.Close()

	return queriedAutoruns
}

func deleteAutorunFromDB(autorun []autoruns.AutorunsInfo) {
	for _, autorun := range autorun {
		zap.S().Info("User deleted: ", autorun.Type)

		conn := db.DbConnect()
		stmt, err := conn.Prepare("DELETE FROM currentautoruns WHERE type = ?")
		if err != nil {
			zap.S().Error("Error deleting autorun from database: ", err)
		}

		_, err = stmt.Exec(autorun.Type)
		if err != nil {
			zap.S().Error("Error deleting autorun from database: ", err)
		}

		defer stmt.Close()
	}

}

func newAutorunIncident(newAutorunIncident []autoruns.AutorunsInfo) {

	payload := AutorunsMonitorPayload{
		NewAutorunsInfo:  autoruns.GetAutoruns(),
		LastLoggedOnUser: usermanagement.GetLastLoggenOnUser(),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		zap.S().Error("Error marshalling payload: ", err)
	}

	newIncident := serialscripter.Incident{
		Name:        "New Autorun Created",
		CurrentTime: time.Now().String(),
		User:        "New Autorun detected",
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
