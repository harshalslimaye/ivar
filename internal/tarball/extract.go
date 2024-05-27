package tarball

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/filesystem"
	"github.com/harshalslimaye/ivar/internal/graph"
	"github.com/mholt/archiver/v3"
)

func ExtractTarball(n *graph.Node, dir string) error {
	packagePath := filepath.Join(n.TargetPath(dir), "package")

	if err := os.MkdirAll(n.TargetPath(dir), 0755); err != nil {
		return fmt.Errorf("unable to create extraction path %s: %s", n.TargetPath(dir), err.Error())
	}

	if err := archiver.Unarchive(n.SourcePath(dir), n.TargetPath(dir)); err != nil {
		return fmt.Errorf("unable to extract %s: %s", n.SourcePath(dir), err.Error())
	}

	if err := filesystem.MoveContents(packagePath, n.TargetPath(dir)); err != nil {
		return fmt.Errorf("unable to move %s to %s: %s", packagePath, n.TargetPath(dir), err.Error())
	}

	if err := os.Remove(packagePath); err != nil {
		return fmt.Errorf("unable to remove %s: %s", packagePath, err.Error())
	}

	return nil
}
