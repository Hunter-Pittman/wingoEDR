# WingoEDR

![alt text](https://github.com/Hunter-Pittman/wingoEDR/blob/main/readme_images/Pasted%20image%2020230328154739.png?raw=true)

## What is WingoEDR
WingoEDR is a multifaceted tool used as both a utility to get key Windows information in a quick and easy to read format, and as a continuous monitoring agent that integrates with another homegrown project, Serial Scripter, all written in Go. The features of WingoEDR are outlined below in multiple sections.
1. Modes
2. Continuous Monitoring
3. Road Map

## Help
This is the help section WingoEDR spits out.
```
Usage of C:\Users\Hunter Pittman\Documents\repos\wingoEDR\wingoEDR.exe:
  -backupdir string
        Enter the path where your backups are going to be stored. (default "C:\\backups")
  -backupitem string
        Enter the path to the file or directory you wish to backup.
  -config string
        Provide path to the config file (default "C:\\Users\\Hunter Pittman\\Documents\\repos\\wingoEDR\\config.json")
  -decompressitem string
        Enter the path to the file or directory you wish to decompress
  -from string
        Enter the start timestamp in the format of YYYY-MM-DDTHH:MM:SS
  -json
        Enter true to output in json format
  -mode string
        List what mode you would like wingoEDR to execute in. The default is to enable continous monitoring. (default "default")
  -offline
        Use this flag to diasble posting to SerialScripter.
  -to string
        Enter the end timestamp in the format of YYYY-MM-DDTHH:MM:SS
```

## Pre-execution
There are several things to know about how WingoEDR operates before running the program. WingoEDR uses/generates 4 other files when it is ran:
1. wingo.db (This is a database containing information regarding the system its on/where it stores monitoring infomation)
2. config.json (This is a config containing necessary context information used by WingoEDR to forward logs, send alerts, and find external resources)
3. externalresources (This is a package containing additional resources not generated by WingoEDR which it downloads from the releases page on the github)
4. logs (The logs directory is generated to contain the output of WingoEDR for either additional analysis later or for archival purposes)

The only file that may need modification is the config.json file for items such as:
* Serial Scripter URL
* Kaspersky Secret
* ELK Stack URL
* Lists

## Modes
WingoEDR can be run in modes each of which run the program in a standalone state that either returns some type of information or causes an action to happen on the system WingoEDR is being run on. To use a mode the syntax is as follows:
```
./wingoEDR.exe --mode=[Name of mode]
```

### Backup
Calling WingoEDR in backup mode has several parameters:
```
./wingoEDR.exe --mode=backup --backupdir=[Path to store area] --backupitem=[File to backup]
```

This mode is relatively straightforward, in short it creates a copy of a chosen file and stores it in the backup directory.
1. Target file exists
2. File is copied and then encrypted
3. File is stored in the backups directory pending restoration

To decompress a file you can use the `decompress` mode.

### Decompress
Calling WingoEDR in decompress mode is straightforward:
```
./wingoEDR.exe --mode=decompress --decompressitem=[Path to compressed file]
```

As was explained in the backup mode section, `decompress` is used to retrieve information from the backed up files. When a file is decompressed it is placed in the directory where you called WingoEDR from.

You can tell when a file has been compressed by WingoEDR when it follow the following naming scheme `compressed_[NAME OF ORIGINAL FILE].[EXTENSION OF ORIGINAL FILE]`. WingoEDR utilized zstandard for compression.

### Chainsaw (Unstable)
`Note: This mode is still under development. The output is still changing, and the time conversion maybe changing as well.`

WingoEDR utilizes Chainsaw, a handy utility for parsing Windows Logs utilizing Sigma detection rules and the mappings created by the Chainsaw team. It's implementation in WingoEDR is basically a wrapper for the Chainsaw binary (stored in the externalresources folder WingoEDR downloads). 

To use the chainsaw mode refer to the following example:
```
./wingoEDR.exe --mode=chainsaw
```
In the above example, Chainsaw is run with the rules categorized as "Bad" (working name...). These rules are derived from the Sigma rule repository and will be augmented with other rules in the future. The rules used for processing can be found in the `./externalresources/rules/Bad` directory.

You also have the option of specifying a time range like so:
```
./wingoEDR.exe --mode=chainsaw --from "2023-03-04T00:00:00" --to "2023-03-05T23:59:59"
```
Use the format specified in the example. Additionally make sure to use your own timezone as WingoEDR will convert to UTC for use with Chainsaw. Time range is a little buggy so double check the output for now until its fixed.

Finally you can get output in a JSON format if you simply add the `--json` flag to the end of the command.
```
./wingoEDR.exe --mode=chainsaw --json
```
This works for the base command as well as the time ranged one.

Example output:
![alt text](https://github.com/Hunter-Pittman/wingoEDR/blob/main/readme_images/Pasted%20image%2020230328180634.png?raw=true)


### Userenum
`Note: This mode needs a strictrer validation (espeically in regards to Acive Directory users).`
This mode returns a list of users, including Active Directory users, if it is run on a Domain Controller. 

Running this one is simple:
```
./wingoEDR.exe --mode=userenum
```

Example Output:
![alt text](https://github.com/Hunter-Pittman/wingoEDR/blob/main/readme_images/Pasted%20image%2020230328181506.png?raw=true)

### Software
This mode returns the installed software for a system.

To run it:
```
./wingoEDR.exe --mode=software
```

Example output:
![alt text](https://github.com/Hunter-Pittman/wingoEDR/blob/main/readme_images/Pasted%20image%2020230328181737.png?raw=true)

## Continuous Monitoring
Continuous monitoring is the other way WingoEDR can be run. This mode has several functions.
1. It reports data back to an associated Serial Scripter web server
2. Makes a baseline of the system it is installed on and reports abnormalities
3. Automatically remediates certain events

When trying to run WingoEDR in continuous monitoring mode there two ways to run it. One simply run the executable `./wingoEDR.exe` and the program will default to continuous monitoring mode or two run WingoEDR as a service (install PowerShell below).

```
function Setup-WingoEDR {

[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

New-Item -ItemType Directory -Path "C:\Program Files\wingoEDR" -Force

$url = "https://github.com/Hunter-Pittman/wingoEDR/releases/download/v0.1.3-alpha/wingoEDR.exe"

$file = "C:\Program Files\wingoEDR\wingoEDR.exe"

$webclient = New-Object System.Net.WebClient

$webclient.DownloadFile($url, $file)

#Expand-Archive -Path $file2 -DestinationPath "C:\Program Files\wingoEDR" -Force

}

function Setup-Service {

New-Service -Name "wingoEDR" -BinaryPathName "C:\Program Files\wingoEDR\wingoEDR.exe" -StartupType Automatic

$process = Start-Process -FilePath "C:\Program Files\wingoEDR\wingoEDR.exe" -NoNewWindow

Start-Sleep -Seconds 10

if(!$process.HasExited){

Stop-Process -Name wingoEDR

}

Start-Service -Name "wingoEDR"

}
```

When WingoEDR is run as a service its auto generated files are found in System32, however config.json can still be edited in the `C:\Program Files\wingoEDR` directory and will be used for running of the service. This will likely change in future, as this split between resources is not ideal.

WingoEDR has a few main components of continuous monitoring listed below:
1. Inventory (This returns essential information about the system to Serial Scripter as part of an effort to contextualize a system with key information)
	1. Host Name
	2. IP
	3. OS
	4. Services
	5. Autoruns (Experimental, not fully implemented) (checks ASEP locations)
	6. Firewall
	7. Shares
	8. Users 
	9. Processes
2. Monitors (These are things that compare a baseline state of a certain component and returns changes)
	1. Autoruns (checks ASEP locations)
	2. Chainsaw (Experimental, not done testing)
	3. Process Monitor (Experimental, not done implementing)
	4. SMB Monitor
	5. Software Monitor (Experimental, not done testing)
	6. User Monitor

All of the modules listed report back to Serial Scripter, but they also write to the WingoEDR log and console (if run in one). Details about what each of these modules return can be found in code and will not be detailed here. Alternatively, run the program yourself and find out!

## Future Road map
This is only the first stage in a long development cycle. The full road map below should explain future improvements.

Phase 1: (WE ARE HERE)
* Implement key OS modules and libraries
	* Network Monitoring 🛑
	* Firewall ✅
	* Shares ✅
	* Processes ✅
	* Registry ✅
	* User Management ⏳
		* User retrieval  ✅
		* Account Manipulation ✅
		* Sessions ⏳
	* Service Management ✅
	* System Health ✅
	* Autoruns ⏳
* Implement key monitors
	* Autoruns Monitor⏳
	* Chainsaw Monitor⏳
	* SMB Monitor ✅
	* Software Monitor ⏳
	* User Monitor ✅
	* Process Monitor ⏳
	* System Health Monitor 🛑
* Implement global logging ✅
* Create a modes parameter ✅
* Serial Scripter API ✅
* Implement Windows service integration ⏳ (Needs param switch and code cleanup) 
* Integrate with Kaspersky ✅
* Integrate with Chainsaw ⏳ 
	* Non time ranged analysis ✅
	* Time ranged analysis ⏳
	* Output refinement ⏳
* Integrate with ELK ⏳ (Needs testing) 
* Utilize an application DB for persistent and programmatic storage ✅
* Create an auto generated configuration system ✅

Phase 2:
* Create an API for the agent (for use with continuous monitoring and Serial Scripter)
* Integrate with Yara-Storm (separate project for remote Yara rule processing)
* Utilize Windows API/syscalls for more operations (phase out as many third party Windows golang libraries as possible)
* Implement better malicious activity algorithms
* Improve action based chainsaw response (Sigma rule trigger response)

Phase 3:
* Improve malicious activity response routes
* Add interactive prompts to certain modes as appropriate (Interactive menus and/or table pagination)





# References
Chainsaw (https://github.com/WithSecureLabs/chainsaw)
Sigma Rule Repository
Yara-Forensics (https://github.com/Xumeiquer/yara-forensics)
