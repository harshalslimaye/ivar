package helper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/packagejson"
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

func GetFileName() string {
	return "package.json"
}

func GetPackageJsonPath() string {
	return filepath.Join(GetCurrentDirPath(), GetFileName())
}

func Exists(path string) bool {
	_, exists := os.Stat(path)

	return !os.IsNotExist(exists)
}

func SameVersionExists(path string, version string) bool {
	pkgJson := packagejson.ReadPackageJson(filepath.Join(path, "package.json"))

	return pkgJson.Version == version
}

func HomeDir() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, ".ivar", "cache")
	}

	return ""
}

func HasHomeDir() bool {
	homeDir := HomeDir()

	if homeDir == "" {
		return false
	}

	_, err := os.Stat(homeDir)

	return !os.IsNotExist(err)
}
