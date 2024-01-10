package utils

import (
	"hash/crc32"
	"os"
	"path/filepath"
)

func IsPathExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func MkdirAll(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Chmod(dirPath, 0777)
	if err != nil {
		return err
	}
	return nil
}

func HashCode(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

func IsDir(filePath string) (bool, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func CreatDir(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Chmod(dirPath, 0777)
	if err != nil {
		return err
	}
	return nil
}

func NotExistToMkdir(dirPath string) error {
	if !IsPathExist(dirPath) {
		return CreatDir(dirPath)
	}
	return nil
}

func RemoveAll(dirPath string) error {
	return os.RemoveAll(dirPath)
}

func CopyDir(source string, destination string) error {
	err := os.MkdirAll(destination, 0755)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, file := range files {
		sourcePath := filepath.Join(source, file.Name())
		destinationPath := filepath.Join(destination, file.Name())

		if file.IsDir() {
			err = CopyDir(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		} else {

			if file.Type()&os.ModeSymlink != 0 {
				link, err := os.Readlink(sourcePath)
				if err != nil {
					return err
				}

				err = os.Symlink(link, destinationPath)
				if err != nil {
					return err
				}
				continue
			}

			data, err := os.ReadFile(sourcePath)
			if err != nil {
				return err
			}

			err = os.WriteFile(destinationPath, data, 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
