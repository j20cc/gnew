package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

type fs struct {
	url, dst string
}

var (
	ds []fs
)

func main() {
	app := &cli.App{
		Name:   "gnew",
		Usage:  "creat new go project quickly by gnew <name>",
		Action: newProject,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func newProject(c *cli.Context) error {
	name := c.Args().Get(0)
	// check arg
	if name == "" {
		return errors.New("project name required")
	}

	// check directory
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return errors.New("project folder existed!")
	}

	initDonwloads(name)

	// mkdir, down stub files
	if err := os.MkdirAll(name+"/cmd/"+name, 0777); err != nil {
		return err
	}
	for _, v := range ds {
		if err := download(v.url, v.dst); err != nil {
			return err
		}
	}

	// replace code

	return nil
}

func initDonwloads(folder string) {
	baseUrl := "https://cdn.jsdelivr.net/gh/lukedever/gnew@master/stubs/"
	d1 := fs{baseUrl + "Makefile", folder + "/Makefile"}
	d2 := fs{baseUrl + "README.md", folder + "/README.md"}
	d3 := fs{baseUrl + "main.stub", fmt.Sprintf("%s/cmd/%s/main.go", folder, folder)}
	ds = append(ds, d1, d2, d3)
}

func download(src, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("download stub file failed")
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
