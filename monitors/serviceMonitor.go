package monitors

import (
	"encoding/json"
	"os"
	"time"
	"wingoEDR/apis/serialscripter"
	"wingoEDR/db"
	"wingoEDR/servicemanager"
	"wingoEDR/usermanagement"

	"go.uber.org/zap"
)

// Payload fields:
// Last logged on user

type ServicePayload struct {
	NewServiceInfo   []servicemanager.WindowsService
	LastLoggedOnUser string
}

func InitServices() {
	services := servicemanager.Servicelister()

	if db.CountTableRecords("currentservices") == 0 {
		servicesToDB(services, false)
	}

}

func MonitorServices() {
	services := servicemanager.Servicelister()
	queriedServices := servicesFromDB()

	if len(services) > len(queriedServices) {
		for _, queriedService := range queriedServices {
			for i, service := range services {
				if service.SCname == queriedService.SCname {
					services[i] = services[len(services)-1]
					services = services[:len(services)-1]
				}
			}
		}
		for _, service := range services {
			zap.S().Info("New service detected: ", service.SCname)
			newServiceIncident([]servicemanager.WindowsService{service})
		}

		servicesToDB(services, false)
	} else if len(services) < len(queriedServices) {
		for _, service := range services {
			for i, queriedService := range queriedServices {
				if service.SCname == queriedService.SCname {
					queriedServices[i] = queriedServices[len(queriedServices)-1]
					queriedServices = queriedServices[:len(queriedServices)-1]
				}
			}
		}

		deleteServiceFromDB(queriedServices)
	} else {
		servicesToDB(services, true)
	}
}

func servicesToDB(services []servicemanager.WindowsService, update bool) {
	var operation string
	if update {
		operation = "UPDATE currentservices SET scname = ?, displayname = ?, statustext = ?, acceptstop = ?, runningpid = ?, port = ?"
	} else {
		operation = "INSERT INTO currentservices (scname, displayname, statustext, acceptstop, runningpid, port) VALUES (?, ?, ?, ?, ?, ?)"
	}

	conn := db.DbConnect()
	stmt, err := conn.Prepare(operation)
	if err != nil {
		zap.S().Error("Error inserting service into database: ", err)
	}

	for _, service := range services {
		_, err := stmt.Exec(service.SCname, service.DisplayName, service.StatusText, service.AcceptStop, service.RunningPID, service.Port)
		if err != nil {
			zap.S().Error("Error inserting service into database: ", err)
		}
	}

	defer stmt.Close()
}

func servicesFromDB() []servicemanager.WindowsService {
	conn := db.DbConnect()
	rows, err := conn.Query("SELECT * FROM currentservices")
	if err != nil {
		zap.S().Error("Error fetching service list from database: ", err)
	}

	defer rows.Close()

	var queriedServices []servicemanager.WindowsService

	for rows.Next() {

		var serviceStructure servicemanager.WindowsService

		err = rows.Scan(&serviceStructure.SCname, &serviceStructure.DisplayName, &serviceStructure.StatusText, &serviceStructure.AcceptStop, &serviceStructure.RunningPID, &serviceStructure.Port)

		if err != nil {
			zap.S().Error(err)
		}

		queriedServices = append(queriedServices, serviceStructure)
	}

	return queriedServices
}

func deleteServiceFromDB(service []servicemanager.WindowsService) {
	for _, service := range service {
		zap.S().Info("service deleted: ", service.SCname)

		conn := db.DbConnect()
		stmt, err := conn.Prepare("DELETE FROM currentservices WHERE scname = ?")
		if err != nil {
			zap.S().Error("Error deleting service from database: ", err)
		}

		_, err = stmt.Exec(service.SCname)
		if err != nil {
			zap.S().Error("Error deleting service from database: ", err)
		}

		defer stmt.Close()
	}

}

func newServiceIncident(newServices []servicemanager.WindowsService) {

	payload := ServicePayload{
		NewServiceInfo:   newServices,
		LastLoggedOnUser: usermanagement.GetLastLoggenOnUser(),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		zap.S().Error("Error marshalling payload: ", err)
	}

	newIncident := serialscripter.Incident{
		Name:        "New Service Created",
		CurrentTime: time.Now().String(),
		User:        "New Service detected",
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
