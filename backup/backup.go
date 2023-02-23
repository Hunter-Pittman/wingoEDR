package backup

import (
	"bufio"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wingoEDR/common"

	"go.uber.org/zap"
)

var (
	dirforbackups = "C:\\backups"
	indexfile     = "C:\\backups\\index.safe"
)

func BackFile(storename string, file string /*, mode int*/) {
	OriginFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	CompressedFile, err := os.Create(dirforbackups + storename)
	if err != nil {
		panic(err)
	}

	PointData := bufio.NewReader(OriginFile)
	common.Compress(PointData, CompressedFile)

	defer OriginFile.Close()
	defer CompressedFile.Close()
}

func ExistsInIndex(indexfile string, file string) string {
	strlist := common.OpenFile(indexfile)

	for _, indexstr := range strlist {
		splittysplit := strings.Split(indexstr, "-:-")
		if splittysplit[0] == file {
			println("exact file exists in index")
			return "newback"
		}
	}
	return "new"
}

func OverWriteModifiedFile(OriginalPath string, FileBackup string) {
	//delete original
	//call modified BackFile function
	os.Remove(OriginalPath)
	BytesToCopy, _ := os.Open(FileBackup)
	NewFile, _ := os.Create(OriginalPath)
	if _, err := io.Copy(NewFile, BytesToCopy); err != nil {
		panic(err)
	}
	defer BytesToCopy.Close()
	defer NewFile.Close()
}

func OverWriteBackup(storename string, file string) {
	f := common.OpenFile(indexfile)
	for _, indexstr := range f {
		var m = make(map[int]string)
		splittysplit := strings.Split(indexstr, "-:-")
		//original file path
		m[0] = splittysplit[0]
		//file backup name
		m[1] = splittysplit[2]
		if file == m[0] {
			os.Remove(dirforbackups + m[1])
			BackFile(m[1], file)
		}
	}
}

func BackDir(file string, overwrite bool) {
	fdir, _ := os.ReadDir(file)

	for _, f := range fdir {
		fpath := filepath.Join(file, f.Name())
		CreateRestorePoint(fpath, overwrite)
	}
}

func GenRandomName() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 15)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func CreateRestorePoint(file string, overwrite bool) {
	stats := common.CheckFile(file)
	if stats.Size != 0 {
		/*
			Index file format:
			Simple ->
			fullpath-:-filename w/extension-:-CompressedBackupName-:-LastModTime-:-hash
			Ex:
			/opt/memento/index.safe-:-index.safe-:-ADZOPRJ13SMF.zst-:-2021-01-01 00:00:00-:-9pN02HFtrhT4EGw+SdIECoj0HV8PBLY8qkZjwaKGRvo=
		*/
		//indexstr := strings.Split(file, "/")
		if stats.Hash == "directory" {
			BackDir(file, overwrite)
		} else {
			strsplit := strings.Split(file, "/")
			storename := strsplit[len(strsplit)-1]

			// /etc/passwd-:-passwd.txt-:-some date-:-hash
			backname := GenRandomName() + ".zst"
			indexstr := file + "-:-" + storename + "-:-" + backname + "-:-" + stats.Time + "-:-" + string(stats.Hash) + "\n"
			newindextstr := []byte(indexstr)

			if _, err := os.Stat(indexfile); os.IsNotExist(err) {
				werr := ioutil.WriteFile(indexfile, newindextstr, 0644)
				if werr != nil {
					panic(werr)
				}

				BackFile(backname, file)
			} else {
				checkresult := ExistsInIndex(indexfile, file)

				switch checkresult {
				case "newback":
					if overwrite {
						zap.S().Info("Overwriting backup for file: " + file)
						//println("Overwriting previous backup of :" + file)
						OverWriteBackup(storename, file)
					} else {
						zap.S().Error("Skipping backup for file, overwrite set to n: " + file)
						println("overwrite is set to n, exiting")
						os.Exit(0)
					}
				case "new":
					appendfile, err := os.OpenFile(indexfile, os.O_APPEND|os.O_WRONLY, 0644)
					if err != nil {
						panic(err)
					}
					appendfile.WriteString(indexstr)
					defer appendfile.Close()

					zap.S().Info("File: " + file + " has been backed up")
					//println("BACKING UP FILE: " + file)

					BackFile(backname, file)
					//PostToServ(m)
				}
			}
		}
	} else {
		println("Nothing to backup :(, file is empty")
	}
}
