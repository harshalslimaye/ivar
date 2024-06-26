package tarball

import (
	"fmt"
	"os"

	"github.com/harshalslimaye/ivar/internal/graph"
)

func DownloadAndExtract(n *graph.Node, dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("unable to create %s: %s", dir, err.Error())
	}

	if n.Graph.HasCache {
		if err := InstallFromCache(n, dir); err != nil {
			if err = InstallFromRegistry(n, dir); err != nil {
				return err
			}
		}
	} else {
		if err := InstallFromRegistry(n, dir); err != nil {
			return err
		}
	}

	return nil
}

func InstallFromRegistry(n *graph.Node, dir string) error {
	if err := DownloadTarball(n, dir); err != nil {
		return fmt.Errorf("unable to download the package %s@%s: %s", n.Name(), n.Version(), err.Error())
	}

	if err := ExtractTarball(n, dir); err != nil {
		return fmt.Errorf("unable to extract the package %s@%s: %s", n.Name(), n.Version(), err.Error())
	}

	if err := DeleteTarball(n.SourcePath(dir)); err != nil {
		return fmt.Errorf("unable to delete the tgz file for %s@%s: %s", n.Name(), n.Version(), err.Error())
	}

	if err := ToCache(n, dir); err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
