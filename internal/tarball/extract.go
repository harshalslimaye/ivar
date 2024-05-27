package tarball

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	if err := moveContents(packagePath, n.TargetPath(dir)); err != nil {
		return fmt.Errorf("unable to move %s to %s: %s", packagePath, n.TargetPath(dir), err.Error())
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

		if entry.IsDir() {
			// Create the target directory if it doesn't exist
			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
					return err
				}
			}
			// Move contents recursively
			if err := moveContents(sourcePath, targetPath); err != nil {
				return err
			}
			// Remove the source directory after moving its contents
			if err := os.Remove(sourcePath); err != nil {
				return err
			}
		} else {
			// Move the file
			if err := os.Rename(sourcePath, targetPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyFile(sourceFile, targetFile string) error {
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return err
	}

	// Ensure the target file has the same permissions as the source file
	sourceInfo, err := os.Stat(sourceFile)
	if err != nil {
		return err
	}
	return os.Chmod(targetFile, sourceInfo.Mode())
}

func CopyContents(sourceDir, targetDir string) error {
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		targetPath := filepath.Join(targetDir, entry.Name())

		if entry.IsDir() {
			// Create the target directory if it doesn't exist
			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
					return err
				}
			}
			// Copy contents recursively
			if err := CopyContents(sourcePath, targetPath); err != nil {
				return err
			}
		} else {
			// Copy the file
			if err := CopyFile(sourcePath, targetPath); err != nil {
				return err
			}
		}
	}

	return nil
}
