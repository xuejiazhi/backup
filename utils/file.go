package utils

import (
	"fmt"
	"github.com/spf13/cast"
	"io/ioutil"
	"os"
	"path/filepath"
)

func FileCount(srcPath string) (cnt int) {
	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			cnt++
		}
		return nil
	})
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	return cnt
}

func GetFileList(srcPath string) (fileList []map[string]string) {
	readerInfos, err := ioutil.ReadDir(srcPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, info := range readerInfos {
		if !info.IsDir() {
			fileList = append(fileList, GetFileStat(srcPath+"/"+info.Name()))
		}
	}
	return
}

func GetFileStat(srcPath string) map[string]string {
	fi, err := os.Stat(srcPath)
	ret := make(map[string]string)
	if err == nil {
		ret["name"] = fi.Name()
		ret["size"] = cast.ToString(fi.Size())
		ret["mode"] = fi.Mode().String()
		ret["modetime"] = fi.ModTime().String()
	}
	return ret
}
