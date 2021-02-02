// Package create_project Create at 2021-01-28 15:03
package create_project

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func InitModel(path, name string) error {
	path = filepath.ToSlash(path)
	err := os.MkdirAll(path+"/"+name+"/internal/model", os.ModePerm)
	if nil != err {
		color.Red("创建model文件夹失败:%s", err.Error())
		return err
	}
	return nil
}
