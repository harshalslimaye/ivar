package helper

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetCurrentDirPath() string {
	curDir, err := os.Getwd()

	if err != nil {
		fmt.Println("Unable to retrieve current folder name or path")
		os.Exit(1)
	}

	return curDir
}

func GetCurrentDirName() string {
	curDirName := filepath.Base(GetCurrentDirPath())

	return curDirName
}

func GetPathSeparator() string {
	return string(filepath.Separator)
}

func GetFileName() string {
	return "package.json"
}

func GetPackageJsonPath() string {
	return GetCurrentDirPath() + GetPathSeparator() + GetFileName()
}
