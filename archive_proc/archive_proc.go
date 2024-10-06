package archive_proc

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
)

func ListArchive(path string) ([]string, error) {
	cmd := exec.Command("lsar", path)
	stdout := bytes.NewBuffer(nil)
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	files := strings.Split(stdout.String(), "\n")
	return files[1:], nil
}

func UnpackArchive(path, dest string) error {
	cmd := exec.Command("unar", path, dest)
	err := cmd.Run()
	return err
}

func PackToRar(path, destPath, filename string) error {
	dest := filepath.Join(destPath, filename)
	cmd := exec.Command(
		"rar", "a", "-v2000000000b",
		"-ep1",
		"-m0", "-r", "-y",
		dest, path)
	err := cmd.Run()
	return err
}
