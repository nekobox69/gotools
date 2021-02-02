// Package create_project Create at 2021-01-28 15:11
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
	serverGo = `
// Package internal Create at {{.CreateTime}}
package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	// not need
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	xormcore "xorm.io/core"
)

var server = &Service{}

func Server() *Service {
	return server
}

type Service struct {
	Config *Config
	Engine *xorm.Engine
	Mode   string
}

func NewService(configPath, mode string) error {
	config, err := LoadConfig(configPath, mode)
	if err != nil {
		Logger.Error(err.Error())
		return err
	}
	server.Config = config

	engine, err := initEngine(config.Db)
	if err != nil {
		Logger.Error(err.Error())
		return err
	}
	server.Engine = engine
	return nil
}

func initEngine(config DBConfig) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(config.DriverName, config.DataSource)
	if nil != err {
		Logger.Error(err.Error())
		return nil, err
	}
	err = engine.Ping()
	if nil != err {
		Logger.Error(err.Error())
		return nil, err
	}
	engine.SetMapper(xormcore.GonicMapper{})
	engine.SetMaxIdleConns(config.MaxIdleCon)
	engine.SetMaxOpenConns(config.MaxCon)
	engine.ShowSQL(true)
	engine.ShowExecTime(true)

	return engine, nil
}

func GetPort() int {
	return server.Config.Port
}
`
	logGo = `
// Package internal Create at {{.CreateTime}}
package internal

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger commom log
var Logger = logrus.StandardLogger()

// ContextHook for log the call context
type ContextHook struct {
	Field  string
	Skip   int
	levels []logrus.Level
}

// NewContextHook use to make an hook
func NewContextHook(levels ...logrus.Level) logrus.Hook {
	hook := ContextHook{
		Field:  "source",
		Skip:   5,
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels implement levels
func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implement fire
func (hook ContextHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}
`
	configGo = `
// Package internal Create at {{.CreateTime}}
package internal

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

// Config server config
type Config struct {
	LogLevel string   {{.Split}}yaml:"log_level"{{.Split}}
	Db       DBConfig {{.Split}}yaml:"db"{{.Split}}
	Port     int      {{.Split}}yaml:"port"{{.Split}}
	Mode     string
}

// DBConfig config of db
type DBConfig struct {
	DataSource string {{.Split}}yaml:"data_source"{{.Split}}
	MaxCon     int    {{.Split}}yaml:"max_con"{{.Split}}
	MaxIdleCon int    {{.Split}}yaml:"max_idle_con"{{.Split}}
	DriverName string {{.Split}}yaml:"driver_name"{{.Split}}
}

// LoadConfig load config from json file
func LoadConfig(path, mode string) (*Config, error) {
	data, err := ioutil.ReadFile(path + mode + ".yaml")

	if err != nil {
		Logger.Error(err.Error())
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		Logger.Error(err.Error())
		return nil, err
	}
	Logger.Infof("%+v", config)
	return &config, nil
}
`
	routerGo = `
// Package internal Create at {{.CreateTime}}
package router

import (
	"{{.Pkg}}/internal"

	"github.com/go-xorm/xorm"
	"github.com/nekobox69/zephyr"
	"github.com/sirupsen/logrus"
)

var Router = zephyr.NewZephyr(internal.Logger)

func init() {
	
}

func Engine() *xorm.Engine {
	return internal.Server().Engine
}

func Logger() *logrus.Logger {
	return internal.Logger
}
`
)

func InitServer(path, name, pkg string) error {
	err := createServer(path, name)
	if nil != err {
		return err
	}
	err = createConfig(path, name)
	if nil != err {
		return err
	}
	err = createLog(path, name)
	if nil != err {
		return err
	}
	err = createRouter(path, name, pkg)
	if nil != err {
		return err
	}
	return nil
}

func createServer(path, name string) error {
	path = filepath.ToSlash(path)
	fileName := path + "/" + name + "/internal/server.go"
	f, err := os.Create(fileName)
	if err != nil {
		color.Red("创建server.go失败:%s", err.Error())
		return err
	}
	defer f.Close()
	t, err := template.New("server.go").Parse(serverGo)
	if err != nil {
		color.Red("创建server.go模板失败:%s", err.Error())
		return err
	}
	err = t.Execute(f, map[string]interface{}{
		"CreateTime": time.Now().Format(timeLayout),
	})
	if err != nil {
		color.Red("生产server.go文件失败:%s", err.Error())
		return err
	}
	cmd := exec.Command("goimports", "-w", fileName)
	err = cmd.Run()
	if nil != err {
		color.Yellow("格式化文件%s失败", fileName)
	}
	return nil
}

func createLog(path, name string) error {
	path = filepath.ToSlash(path)
	fileName := path + "/" + name + "/internal/log.go"
	f, err := os.Create(fileName)

	if err != nil {
		color.Red("创建log.go失败:%s", err.Error())
		return err
	}
	defer f.Close()
	t, err := template.New("log.go").Parse(logGo)
	if err != nil {
		color.Red("创建log.go模板失败:%s", err.Error())
		return err
	}
	err = t.Execute(f, map[string]interface{}{
		"CreateTime": time.Now().Format(timeLayout),
	})
	if err != nil {
		color.Red("生产log.go文件失败:%s", err.Error())
		return err
	}
	cmd := exec.Command("goimports", "-w", fileName)
	err = cmd.Run()
	if nil != err {
		color.Yellow("格式化文件%s失败", fileName)
	}
	return nil
}

func createConfig(path, name string) error {
	path = filepath.ToSlash(path)
	fileName := path + "/" + name + "/internal/config.go"
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		color.Red("创建config.go失败:%s", err.Error())
		return err
	}
	t, err := template.New("config.go").Parse(configGo)
	if err != nil {
		color.Red("创建config.go模板失败:%s", err.Error())
		return err
	}
	err = t.Execute(f, map[string]interface{}{
		"Split":      "`",
		"CreateTime": time.Now().Format(timeLayout),
	})
	if err != nil {
		color.Red("生产config.go文件失败:%s", err.Error())
		return err
	}
	cmd := exec.Command("goimports", "-w", fileName)
	err = cmd.Run()
	if nil != err {
		color.Yellow("格式化文件%s失败", fileName)
	}
	return nil
}

func createRouter(path, name, pkg string) error {
	path = filepath.ToSlash(path)
	err := os.MkdirAll(path+"/"+name+"/internal/router", os.ModePerm)
	if nil != err {
		color.Red("创建router文件夹失败:%s", err.Error())
		return err
	}
	fileName := path + "/" + name + "/internal/router/router.go"
	f, err := os.Create(fileName)

	if err != nil {
		color.Red("创建router.go失败:%s", err.Error())
		return err
	}
	defer f.Close()
	t, err := template.New("router.go").Parse(routerGo)
	if err != nil {
		color.Red("创建router.go模板失败:%s", err.Error())
		return err
	}
	err = t.Execute(f, map[string]interface{}{
		"Pkg":        pkg,
		"CreateTime": time.Now().Format(timeLayout),
	})
	if err != nil {
		color.Red("生产router.go文件失败:%s", err.Error())
		return err
	}
	cmd := exec.Command("goimports", "-w", fileName)
	err = cmd.Run()
	if nil != err {
		color.Yellow("格式化文件%s失败", fileName)
	}
	return nil
}
