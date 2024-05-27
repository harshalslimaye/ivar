package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/Masterminds/semver/v3"
)

var store sync.Map

func GetVersion(name string, version string) string {
	key := fmt.Sprintf("%s@%s", name, version)

	if storedVersion, found := store.Load(key); found {
		if value, okay := storedVersion.(string); okay {
			return value
		}
	}

	exactVersion := FindExactVersion(version, GetAvailableVersions(name))

	store.Store(key, exactVersion)

	return exactVersion
}

func FindExactVersion(constraint string, versions []*semver.Version) string {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		fmt.Printf("Error parsing version constraint %s: %v\n", constraint, err)
		return ""
	}

	var compatibleVersions []*semver.Version
	for _, v := range versions {
		if c.Check(v) {
			compatibleVersions = append(compatibleVersions, v)
		}
	}

	if len(compatibleVersions) == 0 {
		// No compatible versions found
		return ""
	}

	// Sort compatible versions
	sort.Sort(semver.Collection(compatibleVersions))

	// Return the highest compatible version
	return compatibleVersions[len(compatibleVersions)-1].String()
}

func GetAvailableVersions(name string) []*semver.Version {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", name)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code: %d\n", response.StatusCode)
		return nil
	}

	var data struct {
		Versions map[string]struct{} `json:"versions"`
	}

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	// Extract version keys
	var versions []*semver.Version
	for key := range data.Versions {
		v, err := semver.NewVersion(key)
		if err != nil {
			fmt.Printf("Error parsing version %s: %v\n", key, err)
			continue
		}
		versions = append(versions, v)
	}

	return versions
}
