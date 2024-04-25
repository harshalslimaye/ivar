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
			gph := graph.Graph{
				Dependencies: []*graph.Dependency{},
				List:         make(map[string][]*graph.Version),
			}
			graph.BuildGraph(&gph, pkg.Dependencies)

			for k := range gph.List {
				basePath := filepath.Join("node_modules", k)
				if len(gph.List[k]) == 1 {
					DownloadDependency(k, gph.List[k][0].Number, basePath)
				} else {
					for i, value := range gph.List[k] {
						if i == 0 {
							DownloadDependency(k, value.Number, basePath)
						} else {
							for _, dv := range value.Dependents {
								path := filepath.Join("node_modules", dv.Name, "node_modules", k)
								fmt.Println(path)
								DownloadDependency(k, value.Number, filepath.Join(path))
							}
						}
					}
				}
			}
		},
	}

	return cmd
}

func DownloadDependency(name string, version string, downloadDir string) error {
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
