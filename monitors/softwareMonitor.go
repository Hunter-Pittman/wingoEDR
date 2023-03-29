package monitors

import (
	"log"
	"wingoEDR/db"
	"wingoEDR/registrycapture"

	"go.uber.org/zap"
)

func InitSoftware() {
	if db.CountTableRecords("currentsoftware") == 0 {
		softwareList := registrycapture.GetSoftwareSubkeys(`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\Current Version\Uninstall\`)
		softwareToDB(softwareList, false)
	}
}

func softwareToDB(softwareList []registrycapture.InstalledSoftware, update bool) {
	var operation string
	if update {
		operation = "UPDATE currentsoftware SET name = ?, version = ?, installpath = ?, publisher = ?, uninstallstring = ?"
	} else {
		operation = "INSERT INTO currentsoftware (name, version, installpath, publisher, uninstallstring) VALUES (?, ?, ?, ?, ?)"
	}

	conn := db.DbConnect()
	stmt, err := conn.Prepare(operation)
	if err != nil {
		zap.S().Error("Error inserting software list into database: ", err)
	}

	for _, software := range softwareList {
		_, err := stmt.Exec(software.Name, software.Version, software.InstallPath, software.Publisher, software.UninstallString)
		if err != nil {
			zap.S().Error("Error inserting software into database: ", err)
		}
	}

	defer stmt.Close()
}

func deleteSoftwareFromDB(softwareList []registrycapture.InstalledSoftware) {
	for _, software := range softwareList {
		zap.S().Info("Software deleted from database: ", software.Name)

		conn := db.DbConnect()

		stmt, err := conn.Prepare("DELETE FROM currentsoftware WHERE name = ?")
		if err != nil {
			zap.S().Error("Error deleting software from database: ", err)
		}

		_, err = stmt.Exec(software.Name)
		if err != nil {
			zap.S().Error("Error deleting software from database: ", err)
		}

		defer stmt.Close()
	}
}

func softwareFromDB() []registrycapture.InstalledSoftware {
	conn := db.DbConnect()
	rows, err := conn.Query("SELECT * FROM currentsoftware")
	if err != nil {
		zap.S().Error("Error fetching software list from database: ", err)
	}

	defer rows.Close()

	var queriedSoftware []registrycapture.InstalledSoftware

	for rows.Next() {
		var softwareStruct registrycapture.InstalledSoftware
		err = rows.Scan(&softwareStruct.Name, &softwareStruct.Version, &softwareStruct.InstallPath, &softwareStruct.Publisher, &softwareStruct.UninstallString)

		if err != nil {
			log.Fatal(err)
		}

		queriedSoftware = append(queriedSoftware, softwareStruct)
	}
	return queriedSoftware
}

func SoftwareMonitor() {
	softwareList := registrycapture.GetSoftwareSubkeys(`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\Current Version\Uninstall`)
	queriedSoftwareList := softwareFromDB()

	if len(softwareList) > len(queriedSoftwareList) {
		for _, queriedSoftware := range queriedSoftwareList {
			for i, software := range softwareList {
				if software.Name == queriedSoftware.Name {
					softwareList[i] = softwareList[len(softwareList)-1]
					softwareList = softwareList[:len(softwareList)-1]
				}
			}
		}
		for _, software := range softwareList {
			zap.S().Info("New software installed: ", software.Name)
		}
	} else if len(softwareList) < len(queriedSoftwareList) {
		// software deletion
		for _, software := range softwareList {
			for i, queriedSoftware := range queriedSoftwareList {
				if software.Name == queriedSoftware.Name {
					queriedSoftwareList[i] = queriedSoftwareList[len(queriedSoftwareList)-1]
					queriedSoftwareList = queriedSoftwareList[:len(queriedSoftwareList)-1]
				}
			}
		}
		// deleteSoftwareFromDB(queriedSoftwareList)
	} else {
		softwareToDB(softwareList, true)
	}
}
