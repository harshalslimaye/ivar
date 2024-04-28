package tarball

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Use named return values for better readability and error handling
func DownloadTarball(name string, version string, path string) (err error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/-/%s-%s.tgz", name, name, version)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(filepath.Join(path, fmt.Sprintf("%s-%s.tgz", name, version)))
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
