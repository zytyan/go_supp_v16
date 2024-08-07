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
	cmd := exec.Command("lsar", "-x", path, dest)
	err := cmd.Run()
	return err
}

func PackFilesToRar(files []string, destPath, filename string) error {
	dest := filepath.Join(destPath, filename)
	cmd := exec.Command(
		"rar", "a", "-v2147483648b",
		"-m0", "-r", "-y",
		dest)
	cmd.Args = append(cmd.Args, files...)
	err := cmd.Run()
	return err
}

func PackToRar(path, destPath, filename string) error {
	dest := filepath.Join(destPath, filename)
	cmd := exec.Command(
		"rar", "a", "-v2147483648b",
		"-m0", "-r", "-y",

		dest, path)
	err := cmd.Run()
	return err
}
