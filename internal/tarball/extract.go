package tarball

import (
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v3"
)

func ExtractTarball(sourcePath, targetPath string) error {
	packagePath := filepath.Join(targetPath, "package")

	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return err
	}

	if err := archiver.Unarchive(sourcePath, targetPath); err != nil {
		return err
	}

	if err := moveContents(packagePath, targetPath); err != nil {
		return err
	}

	if err := os.Remove(packagePath); err != nil {
		return err
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
