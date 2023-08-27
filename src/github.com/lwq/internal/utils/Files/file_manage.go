package files

import (
	"os"
	"path/filepath"
)

// 删除指定文件夹下的所有文件
func DeleteFilesInDirectory(dirPath string) error {
	fileList := []string{}

	// 遍历目录下的文件
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// 删除文件
	for _, file := range fileList {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}

	return nil
}
