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

// func ExtractTarball(sourcePath, targetPath string) error {
// 	r, err := os.Open(sourcePath)
// 	defer r.Close()

// 	if err != nil {
// 		return err
// 	}

// 	uncompressedStream, err := gzip.NewReader(r)
// 	if err != nil {
// 		return err
// 	}

// 	tarReader := tar.NewReader(uncompressedStream)
// 	for {
// 		header, err := tarReader.Next()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}

// 		targetFilePath := filepath.Join(targetPath, strings.TrimPrefix(header.Name, "package/"))

// 		switch header.Typeflag {
// 		case tar.TypeDir:
// 			if err := os.MkdirAll(targetFilePath, 0755); err != nil {
// 				return fmt.Errorf("ExtractTarGz: MkdirAll() failed: %w", err)
// 			}
// 		case tar.TypeReg:
// 			if err := os.MkdirAll(filepath.Dir(targetFilePath), 0755); err != nil {
// 				return fmt.Errorf("ExtractTarGz: MkdirAll() failed: %w", err)
// 			}
// 			outFile, err := os.Create(targetFilePath)
// 			if err != nil {
// 				return fmt.Errorf("ExtractTarGz: Create() failed: %w", err)
// 			}

// 			if _, err := io.Copy(outFile, tarReader); err != nil {
// 				outFile.Close()
// 				return fmt.Errorf("ExtractTarGz: Copy() failed: %w", err)
// 			}

// 			if err := outFile.Close(); err != nil {
// 				return fmt.Errorf("ExtractTarGz: Close() failed: %w", err)
// 			}
// 		default:
// 			return fmt.Errorf("ExtractTarGz: unknown type: %b in %s", header.Typeflag, header.Name)
// 		}
// 	}

// 	return nil
// }
