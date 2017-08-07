package main

import (
	"github.com/hoisie/web"
	"gopkg.in/yaml.v2"
	"fmt"
	"os"
	"io/ioutil"
	"os/exec"
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

var config Config = Config{}

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

	web.Get("/([0-9a-zA-Z]+)", handle)
	web.Run(config.Listen)
}

func handle(ctx *web.Context, name string) {
	hook, ok := config.Hooks[name]

	if !ok {
		ctx.NotFound("Page not found")
		fmt.Println("1")
		return
	}

	if hook.Password != "" && ctx.Params["password"] != hook.Password {
		ctx.NotFound("Page not found")
		fmt.Println("2")
		return
	}

	// var stderr bytes.Buffer

	cmd := exec.Command("sh", "-c", hook.Execute)
	// cmd.Stderr = &stderr
	cmd.Dir = hook.Directory

	output, err := cmd.CombinedOutput()
	if err != nil {
		ctx.Abort(500, err.Error() + "\n\n")
	}

	ctx.Write(output)
}
