package main

import (
	"fmt"
	"timoni-app/config"
	"timoni-app/console"
	"timoni-app/timoni"
)

func main() {
	conf, err := config.Init()
	if err != nil {
		fmt.Println(err)
	}

	t := timoni.New()

	for _, server := range conf.Task {
		t.AddServer(timoni.Server{
			Name:    server.Name,
			Command: server.Command,
			Args:    server.Args,
			Dir:     server.Dir,
		})
	}

	cmd := console.New(t)
	cmd.Init()
	cmd.Run()
}
