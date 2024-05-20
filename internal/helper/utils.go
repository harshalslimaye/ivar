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

func GetPathSeparator() string {
	return string(filepath.Separator)
}

func GetFileName() string {
	return "package.json"
}

func GetPackageJsonPath() string {
	return GetCurrentDirPath() + GetPathSeparator() + GetFileName()
}

func Exists(path string) bool {
	_, exists := os.Stat(path)

	return !os.IsNotExist(exists)
}

func SameVersionExists(path string, version string) bool {
	pkgJson := packagejson.ReadPackageJson(filepath.Join(path, "package.json"))

	return pkgJson.Version == version
}
