package unzip

import (
	"archive/zip"
	"io"
	"os"

	"go.uber.org/zap"
)

func Unzip(target string) {
	zipFile, err := zip.OpenReader(target)
	if err != nil {
		zap.S().Error("Zip file could not be opened: ", err.Error())
	}
	defer zipFile.Close()

	for _, file := range zipFile.File {
		if !file.FileInfo().IsDir() {
			err := extractFile(file)
			if err != nil {
				zap.S().Error("File could not be extracted: ", err.Error())
			}
		} else {
			err := os.MkdirAll(file.Name, 0755)
			if err != nil {
				zap.S().Error("Directory could not be created: ", err.Error())
			}
		}
	}
}

func extractFile(file *zip.File) error {
	outFile, err := os.Create(file.Name)
	if err != nil {
		return err
	}
	defer outFile.Close()

	inFile, err := file.Open()
	if err != nil {
		return err
	}
	defer inFile.Close()

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		return err
	}

	return nil
}
