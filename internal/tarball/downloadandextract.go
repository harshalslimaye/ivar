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

	if err := DownloadTarball(n, dir); err != nil {
		return fmt.Errorf("unable to download the package %s@%s: %s", n.Name(), n.Version(), err.Error())
	}

	if err := ExtractTarball(n, dir); err != nil {
		return fmt.Errorf("unable to extract the package %s@%s: %s", n.Name(), n.Version(), err.Error())
	}

	if err := DeleteTarball(n.SourcePath(dir)); err != nil {
		return fmt.Errorf("unable to delete the tgz file for %s@%s: %s", n.Name(), n.Version(), err.Error())
	}

	return nil
}
