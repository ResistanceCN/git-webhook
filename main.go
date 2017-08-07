package main

import (
	"github.com/hoisie/web"
	"gopkg.in/yaml.v2"
	"os"
	"io/ioutil"
	"os/exec"
	"bytes"
	"time"
)

type Hook struct {
	Directory  string
	Execute    string
	Password   string
}

type Config struct {
	Listen  string
	Hooks   map[string]Hook
}

var config Config

func main() {
	config = Config{
		Hooks: make(map[string]Hook),
	}

	f, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	f.Close()

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		panic(err)
	}

	web.Match("POST|GET", "/([0-9a-zA-Z]+)", handle)
	web.Run(config.Listen)
}

func handle(ctx *web.Context, name string) {
	hook, ok := config.Hooks[name]

	if !ok {
		ctx.NotFound("Page not found")
		return
	}

	if hook.Password != "" && ctx.Params["password"] != hook.Password {
		ctx.NotFound("Page not found")
		return
	}

	var output bytes.Buffer

	cmd := exec.Command("sh", "-c", hook.Execute)
	cmd.Stdout = &output
	cmd.Stderr = &output
	cmd.Dir = hook.Directory

	err := cmd.Start()

	if err != nil {
		ctx.Abort(500, err.Error() + "\n\n")
	}

	timer := time.AfterFunc(30 * time.Second, func() {
		cmd.Process.Kill()
	})

	cmd.Wait()
	timer.Stop()

	ctx.Write(output.Bytes())
}
