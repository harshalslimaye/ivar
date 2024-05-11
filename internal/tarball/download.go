package tarball

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/graph"
)

// Use named return values for better readability and error handling
func DownloadTarball(node *graph.Node, path string) (err error) {
	res, err := http.Get(node.DownloadPath)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(filepath.Join(path, node.FileName))
	if err != nil {
		return err
	}
	defer func() {
		closeErr := out.Close()
		if err == nil {
			err = closeErr
		}
	}()

	_, err = io.Copy(out, res.Body)
	return err
}
