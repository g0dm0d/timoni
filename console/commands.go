package console

import (
	"fmt"
	"os"
	"strings"
)

func (c *Console) help(argv []string) {
	for key, command := range c.commands {
		fmt.Printf("Command: %s\n", key)
		fmt.Printf("%s\n\n", command.desc)
	}
}

func (c *Console) run(argv []string) {
	if len(argv) < 2 {
		fmt.Printf("The command is not entered correctly\nusage: run <server>\n")
		return
	}
	_, err := c.timoni.GetTask(argv[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	go c.timoni.Run(argv[1])
	fmt.Printf("Server %s is run\n", argv[1])
}

func (c *Console) send(argv []string) {
	if len(argv) < 3 {
		fmt.Printf("The command is not entered correctly\nusage: send <server> <command ...>\n")
		return
	}

	err := c.timoni.SendCommand(argv[1], strings.Join(argv[2:], " "))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (c *Console) stop(argv []string) {
	if len(argv) < 1 {
		fmt.Printf("The command is not entered correctly\nusage: stop <server>\n")
		return
	}

	err := c.timoni.SetAutorestart(argv[1], false)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.timoni.SendCommand(argv[1], "stop")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (c *Console) getLogs(argv []string) {
	if len(argv) < 2 {
		fmt.Printf("The command is not entered correctly\nusage: getLogs <server>\n")
		return
	}

	logs, err := c.timoni.GetLog(argv[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(logs)
}

func (c *Console) getServerList(argv []string) {
	list := c.timoni.GetTaskList()
	fmt.Println(strings.Join(list, ", "))
}

func (c *Console) exit(argv []string) {
	fmt.Println("Exiting...")
	os.Exit(0)
}
