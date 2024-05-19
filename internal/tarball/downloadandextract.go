package tarball

import (
	"fmt"
	"os"

	"github.com/harshalslimaye/ivar/internal/graph"
)

func DownloadAndExtract(n *graph.Node) error {
	if err := os.MkdirAll(n.DownloadDir, 0755); err != nil {
		return fmt.Errorf("unable to create %s: %s", n.DownloadDir, err.Error())
	}

	if err := DownloadTarball(n, n.DownloadDir); err != nil {
		return fmt.Errorf("unable to download the package %s@%s: %s", n.Name(), n.Version(), err.Error())
	}

	if err := ExtractTarball(n); err != nil {
		return fmt.Errorf("unable to extract the package %s@%s: %s", n.Name(), n.Version(), err.Error())
	}

	if err := DeleteTarball(n.SourcePath()); err != nil {
		return fmt.Errorf("unable to delete the tgz file for %s@%s: %s", n.Name(), n.Version(), err.Error())
	}

	return nil
}
