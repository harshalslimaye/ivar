package installcmd

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/packagejson"
	"github.com/spf13/cobra"
)

func InstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "This command installs a package along with its dependencies.",
		Run: func(cmd *cobra.Command, args []string) {
			pkg := packagejson.ReadPackageJson()

			for k, v := range pkg.Dependencies {
				DownloadDependency(k, v)
			}
		},
	}

	return cmd
}

// PackageMetadata represents the structure of package metadata
type PackageMetadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	// Add more fields as needed
}

// DownloadDependency downloads a package from the npm registry and extracts it
func DownloadDependency(packageName string, version string) error {
	// Construct the URL for fetching package metadata
	url := fmt.Sprintf("https://registry.npmjs.org/%s/%s", packageName, version)

	// Send GET request to the npm registry API
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download package metadata: %s", resp.Status)
	}

	// Decode the response body (package metadata) into PackageMetadata struct
	var metadata PackageMetadata
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		return err
	}

	// Create a directory for the downloaded package
	downloadDir := filepath.Join("node_modules", packageName)
	err = os.MkdirAll(downloadDir, 0755)
	if err != nil {
		return err
	}

	// Extract the tarball from the response body
	err = extractTarball(resp.Body, downloadDir)
	if err != nil {
		return err
	}

	return nil
}

// extractTarball extracts the contents of a tarball from the given reader to the specified directory
func extractTarball(reader io.Reader, destDir string) error {
	// Create a gzip reader to decompress the tarball file
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader to read the contents of the tarball
	tarReader := tar.NewReader(gzipReader)

	// Iterate over each file in the tarball
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of tarball
			break
		}
		if err != nil {
			return err
		}

		// Determine the path to extract the file to
		targetPath := filepath.Join(destDir, header.Name)

		// Check if the file is a directory
		if header.Typeflag == tar.TypeDir {
			// Create the directory
			err := os.MkdirAll(targetPath, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			continue
		}

		// Create the file
		file, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Write the file contents to the file
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}

	return nil
}
