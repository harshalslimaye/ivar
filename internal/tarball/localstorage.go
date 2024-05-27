package tarball

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/graph"
	"github.com/harshalslimaye/ivar/internal/helper"
)

func LocalStorage(n *graph.Node, dir string) error {
	if n.Graph.LocalStorageActive {
		localPath := filepath.Join(helper.HomeDir(), fmt.Sprintf("%s@%s", n.Name(), n.Version()))
		if err := os.MkdirAll(localPath, 0755); err == nil {
			copyContents(n.TargetPath(dir), localPath)
		} else {
			return fmt.Errorf("unable to cache %s: %s", n.Package.NameAndVersion(), err.Error())
		}
	}

	return nil
}
