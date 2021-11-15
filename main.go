package main

import (
	"flag"
	"os"

	"github.com/go-sdk/lib/app"
	"github.com/go-sdk/lib/log"
)

type Config struct {
	Tasks []Task `json:"tasks"`

	Path string `json:"-"`
	Help bool   `json:"-"`
}

var config = &Config{}

func init() {
	flag.StringVar(&config.Path, "config", "config.json", "config")
	flag.BoolVar(&config.Help, "help", false, "instructions for use")
	flag.Parse()

	if config.Help {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	a := app.New("cronjob")
	defer a.Recover()

	a.Add(StartTask)

	err := a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
