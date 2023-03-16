package monitors

import (
	"wingoEDR/processes"
)

type EnrichedProcess struct {
	ProcessInfo processes.ProcessInfo
	Exehash     string
	kstatus     string
}

// func InitProcesses() {
// 	if db.CountTableRecords("currentprocesses") == 0 {
// 		procs, err := processes.GetAllProcesses()
// 		if err != nil {
// 			zap.S().Error("Error getting processes: ", err)
// 		}
// 		processesToDB(procs, false)
// 	}
// }

// func ProcessMonitor() {
// 	procs, err := processes.GetAllProcesses()
// 	if err != nil {
// 		zap.S().Error("Error getting processes: ", err)
// 	}

// }

// func processesToDB(shares []processes.ProcessInfo, update bool) {
// 	var operation string
// 	if update {
// 		operation = "UPDATE currentshares SET Name = ?, Ppid = ?, Pid = ?, Exe = ?, Cwd = ?, User = ?"
// 	} else {
// 		operation = "INSERT INTO currentshares (netname, remark, path, type, permissions, maxuses, currentuses) VALUES (?, ?, ?, ?, ?, ?, ?)"
// 	}

// 	conn := db.DbConnect()
// 	stmt, err := conn.Prepare(operation)
// 	if err != nil {
// 		zap.S().Error("Error inserting user into database: ", err)
// 	}

// 	for _, share := range shares {
// 		_, err := stmt.Exec(share.NetName, share.Remark, share.Path, share.Type, share.Permissions, share.MaxUses, share.CurrentUses)
// 		if err != nil {
// 			zap.S().Error("Error inserting share into database: ", err)
// 		}
// 	}

// 	defer stmt.Close()
// }
