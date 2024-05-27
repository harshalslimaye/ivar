package tarball

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/filesystem"
	"github.com/harshalslimaye/ivar/internal/graph"
	"github.com/harshalslimaye/ivar/internal/helper"
)

func ToCache(n *graph.Node, dir string) error {
	if n.Graph.HasCache {
		localPath := filepath.Join(helper.HomeDir(), fmt.Sprintf("%s@%s", n.Name(), n.Version()))
		if err := os.MkdirAll(localPath, 0755); err == nil {
			filesystem.CopyContents(n.TargetPath(dir), localPath)
		} else {
			return fmt.Errorf("unable to cache %s: %s", n.Package.NameAndVersion(), err.Error())
		}
	}

	return nil
}

func InstallFromCache(n *graph.Node, dir string) error {
	if !n.Graph.Cache.IsInCache(n.Name(), n.Version()) {
		return fmt.Errorf("package not available in cache: %s@%s", n.Name(), n.Version())
	}

	if err := filesystem.CopyContents(n.Graph.Cache.Path(n.Name(), n.Version()), n.TargetPath(dir)); err != nil {
		return err
	}

	return nil
}
