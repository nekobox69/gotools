// Package create_project Create at 2021-01-28 15:07
package create_project

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

const mod = `module {{.Project}}

go 1.14

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-xorm/xorm v0.7.9
	github.com/nekobox69/pocket v0.0.0-20210202124852-35d3e9068639
	github.com/nekobox69/zephyr v0.0.0-20210204113825-8986fd3e45e3
	github.com/sirupsen/logrus v1.7.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	xorm.io/core v0.7.3
)
`

func InitMod(path, name string) error {
	path = filepath.ToSlash(path)
	fileName := path + "/" + name + "/go.mod"
	f, err := os.Create(fileName)
	if err != nil {
		color.Red("创建go.mod失败:%s", err.Error())
		return err
	}
	defer f.Close()
	t, err := template.New("go.mod").Parse(mod)
	if err != nil {
		color.Red("创建go.mod模板失败:%s", err.Error())
		return err
	}
	err = t.Execute(f, map[string]interface{}{
		"Project": name,
	})
	if err != nil {
		color.Red("生产go.mod文件失败:%s", err.Error())
		return err
	}
	return nil
}
