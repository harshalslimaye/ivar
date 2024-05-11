package tarball

import (
	"os"

	"github.com/harshalslimaye/ivar/internal/graph"
)

func DownloadAndExtract(node *graph.Node, downloadDir, sourcePath, targetPath string) error {
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return err
	}

	if err := DownloadTarball(node, downloadDir); err != nil {
		return err
	}

	if err := ExtractTarball(sourcePath, targetPath); err != nil {
		return err
	}

	if err := DeleteTarball(sourcePath); err != nil {
		return err
	}

	return nil
}
