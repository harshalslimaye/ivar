package installcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/graph"
	"github.com/harshalslimaye/ivar/internal/helper"
	"github.com/harshalslimaye/ivar/internal/packagejson"
	"github.com/harshalslimaye/ivar/internal/tarball"
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
	downloadDir := filepath.Join("node_modules", name)
	err := os.MkdirAll(downloadDir, 0755)
	if err != nil {
		return err
	}

	tarball.DownloadTarball(name, version, downloadDir)

	r, err := os.Open(fmt.Sprintf(downloadDir+helper.GetPathSeparator()+"%s-%s.tgz", name, version))
	if err != nil {
		fmt.Println("error")
	}
	err = tarball.ExtractTarball(r, helper.GetCurrentDirPath()+helper.GetPathSeparator()+downloadDir)
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
