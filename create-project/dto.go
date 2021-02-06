// Package create_project Create at 2021-01-28 14:57
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
	code = `// Package dto Create at {{.CreateTime}}
package dto

import (
	"encoding/json"
	"fmt"
)

// Resp common response
type Resp struct {
	Code
	Data interface{} {{.Split}}json:"data"{{.Split}}
}

// NewResp param to Object
func NewResp(error *Code, data interface{}) Resp {
	return Resp{
		*error,
		data,
	}
}

// Marshal to json []byte
func (resp *Resp) Marshal() []byte {
	b, err := json.Marshal(resp)
	if err != nil {
		return []byte(err.Error())
	}
	return b
}

// Code object
type Code struct {
	Code int    {{.Split}}json:"code"{{.Split}}
	Msg  string {{.Split}}json:"msg"{{.Split}}
}

// NewCode get object
func NewCode(code int, msg string) *Code {
	return &Code{
		Code: code,
		Msg:  msg,
	}
}

// Code to json
func (e *Code) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf({{.Split}}{"code": %d, "msg": "%s"}{{.Split}}, e.Code, e.Msg)
	}
	return string(b)
}

var (
	// OK success
	OK = &Code{Code: 0, Msg: "ok"}
	// UnsupportedMethodErr 不支持的请求
	UnsupportedMethodErr = &Code{Code: 100000, Msg: "不支持的请求"}
	// SysInternalErr 系统内部错误
	SysInternalErr = &Code{Code: 100001, Msg: "系统内部错误"}
	// ParamErr 参数错误
	ParamErr = &Code{Code: 100002, Msg: "请求参数错误"}
	// LoginErr 登录错误
	LoginErr = &Code{Code: 100003, Msg: "用户名或密码错误"}
	// TokenErr token错误
	TokenErr = &Code{Code: 100004, Msg: "登录过期"}
)
`
)

// InitDto 初始化dto目录
func InitDto(path, name string) error {
	path = filepath.ToSlash(path)
	err := os.MkdirAll(path+"/"+name+"/internal/dto", os.ModePerm)
	if nil != err {
		color.Red("创建dto文件夹失败:%s", err.Error())
		return err
	}
	fileName := path + "/" + name + "/internal/dto/code.go"
	f, err := os.Create(fileName)
	if err != nil {
		color.Red("创建code.go失败:%s", err.Error())
		return err
	}
	defer f.Close()
	t, err := template.New("code.go").Parse(code)
	if err != nil {
		color.Red("创建code.go模板失败:%s", err.Error())
		return err
	}
	err = t.Execute(f, map[string]interface{}{
		"Split":      "`",
		"CreateTime": time.Now().Format(timeLayout),
	})
	if err != nil {
		color.Red("生产code.go文件失败:%s", err.Error())
		return err
	}
	cmd := exec.Command("goimports", "-w", fileName)
	err = cmd.Run()
	if nil != err {
		color.Red("格式化文件%s失败", fileName)
	}
	return nil
}
