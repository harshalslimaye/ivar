package tarball

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/graph"
)

// Use named return values for better readability and error handling
func DownloadTarball(node *graph.Node, path string) (err error) {
	if node.Name() == "@babel/parser" {
		fmt.Println(node.DownloadPath)
	}
	res, err := http.Get(node.DownloadPath)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(filepath.Join(path, fmt.Sprintf("%s-%s.tgz", node.Name(), node.Version())))
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
