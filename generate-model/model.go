// Package generate_model Create at 2021-01-29 8:54
package generate_model

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-xorm/xorm"
)

const (
	modelGo = `
// Package model Create at {{.CreateTime}}
package model

import (
	"errors"
	"fmt"
	"json"

	"{{.Pkg}}/internal"

	"github.com/nekobox69/pocket"
	"github.com/go-xorm/xorm"
)

// {{.StructName}} table {{.Table}}
{{.StructDefine}}

// String print string
func ({{.Idx}} *{{.StructName}}) String() string {
	b,err := json.Marshal({{.Idx}})
	if nil != err {
 		return string(b)
	}
	return fmt.Sprintf("%+v", {{.Idx}})
}

// Add 插入单条记录
func ({{.Idx}} *{{.StructName}}) Add(engine *xorm.Engine) error {
	if nil == engine {
		internal.Logger.Error("engine is empty")
		return errors.New("engine is empty")
	}

	if nil == {{.Idx}}.ID {
		id := pocket.GetUUIDStr()
		{{.Idx}}.ID = &id
	}
	if nil == {{.Idx}}.CreatedAt {
		time := pocket.UnixSecond()
		{{.Idx}}.CreatedAt = &time
	}
	if nil == {{.Idx}}.UpdatedAt {
		{{.Idx}}.UpdatedAt = {{.Idx}}.CreatedAt
	}

	str, param, err := pocket.NewSqlBuilder("{{.Table}}", {{.Idx}}).BuildInsertRow()
	if err != nil {
		internal.Logger.Error(err)
		return err
	}
	param = append([]interface{}{str}, param...)
	_, err = engine.Exec(param...)
	if err != nil {
		internal.Logger.Error(err)
		return err
	}
	return nil
}

// AddTx 插入单条记录
func ({{.Idx}} *{{.StructName}}) AddTx(session *xorm.Session) error {
	if nil == session {
		internal.Logger.Error("engine is empty")
		return errors.New("engine is empty")
	}

	if nil == {{.Idx}}.ID {
		id := pocket.GetUUIDStr()
		{{.Idx}}.ID = &id
	}
	if nil == {{.Idx}}.CreatedAt {
		time := pocket.UnixSecond()
		{{.Idx}}.CreatedAt = &time
	}
	if nil == {{.Idx}}.UpdatedAt {
		{{.Idx}}.UpdatedAt = {{.Idx}}.CreatedAt
	}

	str, param, err := pocket.NewSqlBuilder("{{.Table}}", {{.Idx}}).BuildInsertRow()
	if err != nil {
		internal.Logger.Error(err)
		return err
	}
	param = append([]interface{}{str}, param...)
	_, err = session.Exec(param...)
	if err != nil {
		internal.Logger.Error(err)
		return err
	}
	return nil
}

// EditByID 根据id编辑
func ({{.Idx}} *{{.StructName}}) EditByID(engine *xorm.Engine, item *{{.StructName}}) error {
	if nil == engine || nil == item {
		internal.Logger.Error("param is empty")
		return errors.New("param is empty")
	}
	if nil == {{.Idx}}.ID {
		internal.Logger.Error("id is empty")
		return errors.New("id is empty")
	}
	time := pocket.UnixSecond()
	item.UpdatedAt = &time
	param, err := pocket.XormUpdateParam(item)
	if nil != err {
		internal.Logger.Error(err)
		return err
	}
	if nil == param {
		internal.Logger.Error("no need to update")
		return errors.New("no need to update")
	}
	_, err = engine.Table("{{.Table}}").Where("disabled=0 AND id=?", {{.Idx}}.ID).Update(param)
	if nil != err {
		internal.Logger.Error(err)
		return err
	}
	return nil
}

// EditByIDTx 根据id编辑
func ({{.Idx}} *{{.StructName}}) EditByIDTx(session *xorm.Session, item *{{.StructName}}) error {
	if nil == session || nil == item {
		internal.Logger.Error("param is empty")
		return errors.New("param is empty")
	}
	if nil == {{.Idx}}.ID {
		internal.Logger.Error("id is empty")
		return errors.New("id is empty")
	}
	time := pocket.UnixSecond()
	item.UpdatedAt = &time
	param, err := pocket.XormUpdateParam(item)
	if nil != err {
		internal.Logger.Error(err)
		return err
	}
	if nil == param {
		internal.Logger.Error("no need to update")
		return errors.New("no need to update")
	}

	_, err = session.Table("{{.Table}}").Where("disabled=0 AND id=?", {{.Idx}}.ID).Update(param)
	if nil != err {
		internal.Logger.Error(err)
		return err
	}
	return nil
}

// EditByCustom 根据自定义条件编辑
func ({{.Idx}} {{.StructName}}) EditByCustom(engine *xorm.Engine, item *{{.StructName}}, where string, whereParam []interface{}) error {
	if nil == engine || nil == item {
		internal.Logger.Error("param is empty")
		return errors.New("param is empty")
	}
	if nil == whereParam {
		whereParam = make([]interface{}, 0)
	}

	time := pocket.UnixSecond()
	item.UpdatedAt = &time
	param, err := pocket.XormUpdateParam(item)
	if nil != err {
		internal.Logger.Error(err)
		return err
	}
	if nil == param {
		internal.Logger.Error("no need to update")
		return errors.New("no need to update")
	}
	_, err = engine.Table("{{.Table}}").Where(where, whereParam...).Update(param)
	if nil != err {
		internal.Logger.Error(err)
		return err
	}
	return nil
}

// EditByCustomTx 根据自定义条件编辑
func ({{.Idx}} {{.StructName}}) EditByCustomTx(session *xorm.Session, item *{{.StructName}}, where string, whereParam []interface{}) error {
	if nil == session || nil == item {
		internal.Logger.Error("param is empty")
		return errors.New("param is empty")
	}
	if nil == whereParam {
		whereParam = make([]interface{}, 0)
	}

	time := pocket.UnixSecond()
	item.UpdatedAt = &time
	param, err := pocket.XormUpdateParam(item)
	if nil != err {
		internal.Logger.Error(err)
		return err
	}
	if nil == param {
		internal.Logger.Error("no need to update")
		return errors.New("no need to update")
	}
	_, err = session.Table("{{.Table}}").Where(where, whereParam...).Update(param)
	if nil != err {
		internal.Logger.Error(err)
		return err
	}
	return nil
}

// GetByID 根据id查询
func ({{.Idx}} {{.StructName}}) GetByID(engine *xorm.Engine) (*{{.StructName}}, error) {
	if nil == engine || nil == {{.Idx}}.ID {
		internal.Logger.Error("param is empty")
		return nil, errors.New("param is empty")
	}
	var {{.TableObj}} {{.StructName}}
	b, err := engine.Table("{{.Table}}").Where("id=? AND disabled=0", {{.Idx}}.ID).Get(&{{.TableObj}})
	if nil != err {
		internal.Logger.Error(err)
		return nil, nil
	}
	if !b {
		internal.Logger.Error("not found")
		return nil, nil
	}
	return &{{.TableObj}}, nil
}

// GetList 分页查询列表
func ({{.Idx}} {{.StructName}}) GetList(engine *xorm.Engine, where string, whereParam []interface{}, page, pageSize int, sort map[string][]string) ([]{{.StructName}}, int64, error) {
	list := new([]{{.StructName}})
	session := engine.Table("{{.Table}}").Where(where, whereParam...).Limit(pageSize, page)
	if asc, ok := sort["asc"]; ok {
		session.Asc(asc...)
	}
	if desc, ok := sort["desc"]; ok {
		session.Desc(desc...)
	}
	count, err := session.FindAndCount(list)
	if nil != err {
		core.Logger.Error(err)
		return nil, count, err
	}

	return *list, count, nil
}
`
	driver     = "mysql"
	dataSource = "%s:%s@tcp(%s)/%s"
	timeLayout = "2006-01-02 15:04:05"
)

type DBTpl struct {
	Package      string
	StructDefine string
	Idx          string
	Table        string
	TableObj     string
	StructName   string
	CreateTime   string
}

func GenerateModel(host, db, user, pwd, pkg, path string) error {
	color.Blue("开始建立数据库连接")
	engine, err := xorm.NewEngine(driver, fmt.Sprintf(dataSource, user, pwd, host, db))
	if nil != err {
		color.Red("初始化engine失败:%s", err.Error())
		return err
	}
	err = engine.Ping()
	if nil != err {
		color.Red("连接数据库失败:%s", err.Error())
		return err
	}
	color.Blue("建立数据库连接成功")
	err = model(engine, pkg, path)
	return err
}

func model(engine *xorm.Engine, pkg, path string) error {
	color.Blue("开始读取表定义")
	tables, err := engine.DBMetas()
	if nil != err {
		color.Red("获取数据库结构失败:", err.Error())
		return err
	}
	for _, v := range tables {
		color.Blue("开始生成表%s对应的model", v.Name)
		fileName := fmt.Sprintf("db_%s.go", v.Name)
		structName := sqlParamToGoParam(v.Name)
		structDefine := fmt.Sprintf("type %s struct {\n", structName)
		structRow := ""
		for _, c := range v.Columns() {
			structRow += fmt.Sprintf("	%s		%s		`json:%s db:%s`\n",
				sqlParamToGoParam(c.Name), sqlTypeToGoType(c.SQLType.Name), fmt.Sprintf("\"%s\"", c.Name), fmt.Sprintf("\"%s\"", c.Name))
		}
		structRow += "}"
		tpl := DBTpl{
			Package:      pkg,
			StructDefine: structDefine + structRow,
			Idx:          v.Name[:1],
			Table:        v.Name,
			TableObj:     fmt.Sprintf("%s%s", strings.ToLower(structName[:1]), structName[1:]),
			StructName:   structName,
			CreateTime:   time.Now().Format(timeLayout),
		}
		err = dbSaveToFile(path+fileName, tpl)
		color.Blue("结束生成表%s对应的model", v.Name)
		if nil != err {
			return err
		}
	}
	return err
}

func dbSaveToFile(fileName string, tpl DBTpl) error {
	f, err := os.Create(fileName)
	if err != nil {
		color.Red("创建文件%s失败：%s", fileName, "", err.Error())
		return err
	}
	defer f.Close()
	t, err := template.New("model.go").Parse(modelGo)
	if err != nil {
		color.Red("创建模板失败:%s", err.Error())
		return err
	}
	err = t.Execute(f, map[string]interface{}{
		"Pkg":          tpl.Package,
		"StructDefine": template.HTML(tpl.StructDefine),
		"Table":        tpl.Table,
		"TableObj":     tpl.TableObj,
		"Idx":          tpl.Idx,
		"StructName":   tpl.StructName,
		"CreateTime":   tpl.CreateTime,
	})
	if err != nil {
		color.Red("生成文件%s失败：%s", fileName, err.Error())
		return err
	}
	cmd := exec.Command("goimports", "-w", fileName)
	err = cmd.Run()
	if nil != err {
		color.Yellow("格式化文件%s失败", fileName)
	}
	return nil
}
