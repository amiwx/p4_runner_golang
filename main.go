package main

import (
	"fmt"
	"os"

	"github.com/amiwx/p4_runner_golang/app"
	"github.com/amiwx/p4_runner_golang/config"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("fatal error config file: config.yml \n", err)
		os.Exit(1)
	}

	a := app.App{}
	a.Initialize(config)
	a.Run(":8000")
}
