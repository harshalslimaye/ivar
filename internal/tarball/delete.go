package tarball

import "os"

// Use named return values for better readability and error handling
func DeleteTarball(sourcePath string) (err error) {
	err = os.Remove(sourcePath)
	return err
}
