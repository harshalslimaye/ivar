package tarball

import (
	"fmt"
	"os"

	"github.com/harshalslimaye/ivar/internal/graph"
)

func DownloadAndExtract(node *graph.Node, downloadDir, sourcePath, targetPath string) error {
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("unable to create %s: %s", downloadDir, err.Error())
	}

	if err := DownloadTarball(node, downloadDir); err != nil {
		return fmt.Errorf("unable to download the package %s@%s: %s", node.Name(), node.Version(), err.Error())
	}

	if err := ExtractTarball(sourcePath, targetPath); err != nil {
		return fmt.Errorf("unable to extract the package %s@%s: %s", node.Name(), node.Version(), err.Error())
	}

	if err := DeleteTarball(sourcePath); err != nil {
		return fmt.Errorf("unable to delete the tgz file for %s@%s: %s", node.Name(), node.Version(), err.Error())
	}

	return nil
}
