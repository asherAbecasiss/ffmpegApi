package main

import (
	"io/ioutil"
	"os"
)

func ensureDir(dirName string) error {

	err := os.Mkdir(dirName, 0777)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
func CountFileInFolder(path string) int {
	files, _ := ioutil.ReadDir(path)

	return len(files)
}
