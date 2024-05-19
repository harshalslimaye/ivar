package tarball

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/graph"
	"github.com/mholt/archiver/v3"
)

func ExtractTarball(n *graph.Node) error {
	packagePath := filepath.Join(n.TargetPath(), "package")

	if err := os.MkdirAll(n.TargetPath(), 0755); err != nil {
		return fmt.Errorf("unable to create extraction path %s: %s", n.TargetPath(), err.Error())
	}

	if err := archiver.Unarchive(n.SourcePath(), n.TargetPath()); err != nil {
		return fmt.Errorf("unable to extract %s: %s", n.SourcePath(), err.Error())
	}

	if err := moveContents(packagePath, n.TargetPath()); err != nil {
		return fmt.Errorf("unable to move %s to %s: %s", packagePath, n.TargetPath(), err.Error())
	}

	if err := os.Remove(packagePath); err != nil {
		return fmt.Errorf("unable to remove %s: %s", packagePath, err.Error())
	}

	return nil
}

func moveContents(sourceDir, targetDir string) error {
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		targetPath := filepath.Join(targetDir, entry.Name())

		if err := os.Rename(sourcePath, targetPath); err != nil {
			return err
		}
	}

	return nil
}
