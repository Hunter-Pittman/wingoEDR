package monitors

import (
	"wingoEDR/db"
	"wingoEDR/software"

	"github.com/blockloop/scan/v2"
	"go.uber.org/zap"
)

func InitSoftware() {
	softwareList := software.GetInstalledSofware()
	if db.CountTableRecords("currentsoftware") == 0 {
		//softwareList := registrycapture.GetSoftwareSubkeys(`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\Current Version\Uninstall`)
		softwareToDB(softwareList, false)
	}
}

func softwareToDB(softwareList []software.Software, update bool) {
	var operation string
	if update {
		operation = "UPDATE currentsoftware SET name = ?, version = ?, vendor = ?"
	} else {
		operation = "INSERT INTO currentsoftware (name, version, vendor) VALUES (?, ?, ?)"
	}

	conn := db.DbConnect()
	stmt, err := conn.Prepare(operation)
	if err != nil {
		zap.S().Error("Error inserting software list into database: ", err)
	}

	for _, software := range softwareList {
		_, err := stmt.Exec(software.Name, software.Version, software.Vendor)
		if err != nil {
			zap.S().Error("Error inserting software into database: ", err)
		}
	}

	defer stmt.Close()
}

func deleteSoftwareFromDB(softwareList []software.Software) {
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

func softwareFromDB() []software.Software {
	conn := db.DbConnect()
	rows, err := conn.Query("SELECT * FROM currentsoftware")
	if err != nil {
		zap.S().Error("Error fetching software list from database: ", err)
	}

	defer rows.Close()

	var queriedSoftware []software.Software
	err1 := scan.Rows(&queriedSoftware, rows)
	if err1 != nil {
		zap.S().Error("Error scanning software list from database")
	}

	/*for rows.Next() {
		var softwareStruct software.Software
		err = rows.Scan(&softwareStruct.Name, &softwareStruct.Version, &softwareStruct.Vendor)

		if err != nil {
			log.Fatal(err)
		}

		queriedSoftware = append(queriedSoftware, softwareStruct)
	}*/
	return queriedSoftware
}

func SoftwareMonitor() {
	//softwareList := registrycapture.GetSoftwareSubkeys(`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\Current Version\Uninstall`)
	softwareList := software.GetInstalledSofware()
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
