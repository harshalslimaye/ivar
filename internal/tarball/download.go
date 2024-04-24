package tarball

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/harshalslimaye/ivar/internal/helper"
)

func DownloadTarball(name string, version string, path string) error {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/-/%s-%s.tgz", name, name, version)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(fmt.Sprintf(path+helper.GetPathSeparator()+"%s-%s.tgz", name, version))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	return nil
}
