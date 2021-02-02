// Package create_project Create at 2021-01-28 14:28
package create_project

import (
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/fatih/color"
)

const (
	timeLayout = "2006-01-02 15:04:05"
	configYaml = `
port: 4000
log_level: "debug"
db:
  data_source: "root:123456@tcp(localhost:3306)/{{.DB}}"
  max_con: 20
  max_idle_con: 5
  driver_name: "mysql"
`
	mainGO = `
// Package main Create at {{.CreateTime}}
package main

import (
	"fmt"
	"net/http"
	"os"

	"{{.Pkg}}/internal"
	"{{.Pkg}}/internal/router"

	"github.com/nekobox69/zephyr"
	"github.com/sirupsen/logrus"
)

var (
	DepMode   string
	BuildTime string
	GoVersion string
	Version   string
)

func init() {
	internal.Logger.SetFormatter(&logrus.JSONFormatter{})
	internal.Logger.AddHook(internal.NewContextHook())
}

func main() {
	DepMode := os.Getenv("DEP_MODE")
	if len(DepMode) == 0 {
		DepMode = internal.StandAlone
	}
	internal.Logger.Info(DepMode)
	err := internal.NewService("./config/", DepMode)
	if nil != err {
		internal.Logger.Fatal(err.Error())
	}

	profile := zephyr.NewProfile("/{{.Project}}/about", Version, GoVersion, BuildTime, DepMode)
	router.Router.SetProfile(profile)
	http.HandleFunc("/", router.Router.ServeHTTP)
	internal.Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", internal.GetPort()), nil))
}
`
)

// InitCmd 初始化cmd目录
func InitCmd(path, name, db, pkg string) error {
	path = filepath.ToSlash(path)
	err := os.MkdirAll(path+"/"+name+"/cmd/config", os.ModePerm)
	if nil != err {
		color.Red("创建cmd文件夹失败:", err.Error())
		return err
	}
	err = createConfigYaml(path, name, db)
	if nil != err {
		return err
	}
	return createMain(path, name, pkg)
}

func createConfigYaml(path, name, db string) error {
	f, err := os.Create(path + "/" + name + "/cmd/config/dev.yaml")
	if err != nil {
		color.Red("创建dev.yaml失败:%s", err.Error())
		return err
	}
	defer f.Close()
	t, err := template.New("dev.yaml").Parse(configYaml)
	if err != nil {
		color.Red("创建dev.yaml模板失败:%s", err.Error())
		return err
	}
	err = t.Execute(f, map[string]interface{}{
		"DB": db,
	})
	if err != nil {
		color.Red("生产dev.yaml文件失败:%s", err.Error())
		return err
	}

	ft, err := os.Create(path + "/" + name + "/cmd/config/test.yaml")
	if nil != err {
		color.Red("创建test.yaml失败：%s", err.Error())
		return err
	}
	defer ft.Close()

	err = t.Execute(ft, map[string]interface{}{
		"DB": db,
	})
	if err != nil {
		color.Red("生产test.yaml文件失败:%s", err.Error())
		return err
	}

	fp, err := os.Create(path + "/" + name + "/cmd/config/prod.yaml")
	if nil != err {
		color.Red("创建prod.yaml失败：%s", err.Error())
		return err
	}
	defer fp.Close()
	err = t.Execute(fp, map[string]interface{}{
		"DB": db,
	})
	if err != nil {
		color.Red("生产prod.yaml文件失败:%s", err.Error())
		return err
	}
	return nil
}

func createMain(path, name, pkg string) error {
	fileName := path + "/" + name + "/cmd/main.go"
	f, err := os.Create(fileName)
	if nil != err {
		color.Red("创建main.go失败：%s", err.Error())
		return err
	}
	defer f.Close()
	t, err := template.New("main.go").Parse(mainGO)
	if err != nil {
		color.Red("创建main.go模板失败:%s", err.Error())
		return err
	}
	err = t.Execute(f, map[string]interface{}{
		"Pkg":        pkg,
		"CreateTime": time.Now().Format(timeLayout),
	})
	if err != nil {
		color.Red("生产main.go文件失败:%s", err.Error())
		return err
	}
	cmd := exec.Command("goimports", "-w", fileName)
	err = cmd.Run()
	if nil != err {
		color.Red("格式化文件%s失败", fileName)
	}
	return nil
}
