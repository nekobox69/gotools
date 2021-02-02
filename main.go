// Package test Create at 2020-12-09 8:44
package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	cp "github.com/nekobox69/gotools/create-project"
	generatemodel "github.com/nekobox69/gotools/generate-model"

	"github.com/fatih/color"
	cli "github.com/urfave/cli/v2"
)

var (
	name string
	path string
	db   string
	pkg  string
	host string
	user string
	pwd  string
)

func main() {
	app := cli.App{
		Name:    "Go Tools",
		Version: "1.0.0",
		Description: `
Go Tools
1.初始化Go Web工程
2.生成model
`,
		Commands: []*cli.Command{
			{
				Name:        "init",
				Description: "初始化Go Web工程",
				Action: func(context *cli.Context) error {
					if "" == pkg {
						pkg = name
					}
					return initProject()
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "name",
						Aliases:     []string{"n"},
						Usage:       "demo",
						Destination: &name,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "path",
						Usage:       "/tmp/",
						Destination: &path,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "package",
						Aliases:     []string{"p"},
						Usage:       "github.com/aa/a",
						Destination: &pkg,
					},
					&cli.StringFlag{
						Name:        "db",
						Aliases:     []string{"d"},
						Usage:       "demo",
						Destination: &db,
					},
				},
			},
			{
				Name:        "generate-model",
				Description: "生成model",
				Action: func(context *cli.Context) error {
					return generatemodel.GenerateModel(host, db, user, pwd, pkg, filepath.ToSlash(path))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "path",
						Usage:       "/tmp/",
						Destination: &path,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "package",
						Aliases:     []string{"p"},
						Usage:       "github.com/aa/a",
						Destination: &pkg,
					},
					&cli.StringFlag{
						Name:        "db",
						Aliases:     []string{"d"},
						Usage:       "demo",
						Destination: &db,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "host",
						Usage:       "127.0.0.1:3306",
						Destination: &host,
						Value:       "127.0.0.1:3306",
					},
					&cli.StringFlag{
						Name:        "user",
						Aliases:     []string{"u"},
						Usage:       "root",
						Destination: &user,
						Value:       "root",
					},
					&cli.StringFlag{
						Name:        "password",
						Aliases:     []string{"pwd"},
						Usage:       "123456",
						Destination: &pwd,
						Value:       "123456",
						Required:    true,
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func initProject() error {
	color.Blue("开始生成cmd")
	err := cp.InitCmd(path, name, db, pkg)
	if nil != err {
		return err
	}
	color.Blue("结束生成cmd")
	color.Blue("开始生成dto")
	err = cp.InitDto(path, name)
	if nil != err {
		return err
	}
	color.Blue("结束生成dto")
	color.Blue("开始生成model")
	err = cp.InitModel(path, name)
	if nil != err {
		return err
	}
	color.Blue("结束生成model")
	color.Blue("开始生成server")
	err = cp.InitServer(path, name, pkg)
	if nil != err {
		return err
	}
	color.Blue("结束生成server")
	color.Blue("开始生成mod")
	err = cp.InitMod(path, name)
	if nil != err {
		return err
	}
	color.Blue("结束生成mod")
	color.Green("创建项目完成")
	return nil
}
