package packagejson

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iancoleman/orderedmap"
)

type PackageJson struct {
	Name            string
	Version         string
	Description     string
	Main            string
	Scripts         map[string]string
	Dependencies    map[string]string
	DevDependencies map[string]string
	Repository      string
	Keywords        []string
	License         string
	Author          string
}

func (pkg *PackageJson) ToInitJson() ([]byte, error) {
	jsonMap := orderedmap.New()
	jsonMap.Set("name", pkg.Name)
	jsonMap.Set("version", pkg.Version)
	jsonMap.Set("description", pkg.Description)
	jsonMap.Set("main", pkg.Main)
	jsonMap.Set("repository", pkg.Repository)

	if pkg.Keywords != nil && len(pkg.Keywords) > 0 {
		jsonMap.Set("keywords", pkg.Keywords)
	}

	jsonMap.Set("author", pkg.Author)
	jsonMap.Set("license", pkg.License)

	return json.MarshalIndent(jsonMap, "", "  ")
}

func (pkg *PackageJson) PrintInitJson() {
	jsonBytes, err := pkg.ToInitJson()
	if err != nil {
		fmt.Print("Error converting to JSON:", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonBytes))
}

func (pkg *PackageJson) WriteToFile(filename string) error {
	jsonData, err := pkg.ToInitJson()
	if err != nil {
		return err
	}

	file, createErr := os.Create(filename)
	if createErr != nil {
		return createErr
	}
	defer file.Close()

	_, writeErr := file.Write(jsonData)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

func GetNewPackageJson(hasDefault bool, dirName string) *PackageJson {
	pkgjson := PackageJson{}

	if hasDefault {
		pkgjson.Name = dirName
		pkgjson.Version = "1.0.0"
		pkgjson.Description = ""
		pkgjson.Main = "index.js"
		pkgjson.Repository = ""
		pkgjson.License = "MIT"
	}

	return &pkgjson
}

func Exists() bool {
	_, err := os.Stat("package.json")
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return false
}

func ReadPackageJson(path string) *PackageJson {
	content, err := os.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	pkg := PackageJson{}
	json.Unmarshal(content, &pkg)

	return &pkg
}
