package tarball

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractTarball(gzipStream io.Reader, targetPath string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetFilePath := filepath.Join(targetPath, strings.TrimPrefix(header.Name, "package/"))

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetFilePath, 0755); err != nil {
				return fmt.Errorf("ExtractTarGz: MkdirAll() failed: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetFilePath), 0755); err != nil {
				return fmt.Errorf("ExtractTarGz: MkdirAll() failed: %w", err)
			}
			outFile, err := os.Create(targetFilePath)
			if err != nil {
				return fmt.Errorf("ExtractTarGz: Create() failed: %w", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("ExtractTarGz: Copy() failed: %w", err)
			}

			if err := outFile.Close(); err != nil {
				return fmt.Errorf("ExtractTarGz: Close() failed: %w", err)
			}
		default:
			return fmt.Errorf("ExtractTarGz: unknown type: %b in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}
