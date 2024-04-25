package vercon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/Masterminds/semver/v3"
)

type Vercon struct {
	cache map[string][]*semver.Version
}

func (v *Vercon) GetVersion(name string, version string) string {
	return v.FindExactVersion(version, v.GetAvailableVersions(name))
}

func (v *Vercon) FindExactVersion(constraint string, versions []*semver.Version) string {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		// Handle invalid constraint
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

func (v *Vercon) GetAvailableVersions(name string) []*semver.Version {
	value, ok := v.cache[name]

	if ok {
		return value
	} else {
		url := fmt.Sprintf("https://registry.npmjs.org/%s", name)
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making HTTP request:", err)
		}
		defer response.Body.Close()

		// Read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
		}

		// Parse the JSON response
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}

		// Extract keys of the "versions" array
		versions := data["versions"].(map[string]interface{})
		keys := make([]*semver.Version, 0, len(versions))
		for key := range versions {
			keys = append(keys, semver.MustParse(key))
		}
		v.cache[name] = keys

		return keys
	}
}

func NewVercon() *Vercon {
	return &Vercon{
		cache: make(map[string][]*semver.Version),
	}
}
