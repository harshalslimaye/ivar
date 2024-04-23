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
	"strings"

	"github.com/harshalslimaye/ivar/internal/graph"
	"github.com/harshalslimaye/ivar/internal/helper"
	"github.com/harshalslimaye/ivar/internal/packagejson"
	"github.com/spf13/cobra"
)

func InstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "This command installs a package along with its dependencies.",
		Run: func(cmd *cobra.Command, args []string) {
			pkg := packagejson.ReadPackageJson()
			list := graph.DependencyList{
				Dependencies: []*graph.Dependency{},
				Visited:      make(map[string]string),
			}
			graph.BuildList(&list, pkg.Dependencies)

			for _, v := range list.Dependencies {
				DownloadDependency(v.Name, v.Version)
			}
		},
	}

	return cmd
}

func DownloadDependency(name string, version string) error {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/%s", name, version)
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	var dep graph.Dependency
	err = json.NewDecoder(res.Body).Decode(&dep)
	if err != nil {
		return err
	}

	downloadDir := filepath.Join("node_modules", name)
	err = os.MkdirAll(downloadDir, 0755)
	if err != nil {
		return err
	}

	DownloadTarball(name, version, downloadDir)

	r, err := os.Open(fmt.Sprintf(downloadDir+helper.GetPathSeparator()+"%s-%s.tgz", name, version))
	if err != nil {
		fmt.Println("error")
	}
	err = ExtractTarball(r, helper.GetCurrentDirPath()+helper.GetPathSeparator()+downloadDir)
	if err != nil {
		return err
	}
	r.Close()

	// Deleting the tarball once installation is complete
	tarballPath := fmt.Sprintf(downloadDir+helper.GetPathSeparator()+"%s-%s.tgz", name, version)
	if err := os.Remove(tarballPath); err != nil {
		return fmt.Errorf("DownloadDependency: Failed to delete tarball: %w", err)
	}

	return nil
}

func DownloadTarball(name string, version string, path string) error {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/-/%s-%s.tgz", name, name, version)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(fmt.Sprintf(path+helper.GetPathSeparator()+"%s-%s.tgz", name, version))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func ExtractTarball(gzipStream io.Reader, targetPath string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetFilePath := filepath.Join(targetPath, strings.TrimPrefix(header.Name, "package/"))

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetFilePath, 0755); err != nil {
				return fmt.Errorf("ExtractTarGz: MkdirAll() failed: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetFilePath), 0755); err != nil {
				return fmt.Errorf("ExtractTarGz: MkdirAll() failed: %w", err)
			}
			outFile, err := os.Create(targetFilePath)
			if err != nil {
				return fmt.Errorf("ExtractTarGz: Create() failed: %w", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("ExtractTarGz: Copy() failed: %w", err)
			}

			if err := outFile.Close(); err != nil {
				return fmt.Errorf("ExtractTarGz: Close() failed: %w", err)
			}
		default:
			return fmt.Errorf("ExtractTarGz: unknown type: %b in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}
