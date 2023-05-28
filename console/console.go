package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"timoni-app/timoni"
)

type Console struct {
	timoni   *timoni.Timoni
	commands map[string]Command
}

type Command struct {
	function func(argv []string)
	desc     string
}

func New(t *timoni.Timoni) *Console {
	return &Console{
		timoni:   t,
		commands: make(map[string]Command),
	}
}

func (c *Console) Init() {
	c.commands["help"] = Command{
		function: c.help,
		desc:     "HELP ME!!!",
	}
	c.commands["run"] = Command{
		function: c.run,
		desc:     "run <server>",
	}
	c.commands["send"] = Command{
		function: c.send,
		desc:     "send <server> <command>\ncommand to send a command to the server\nExample usage: send lobby say hello",
	}
	c.commands["exit"] = Command{
		function: c.exit,
		desc:     "exit from app",
	}
	c.commands["stop"] = Command{
		function: c.stop,
		desc:     "stop <server>\nstopping server",
	}
	c.commands["getlogs"] = Command{
		function: c.getLogs,
		desc:     "getlogs <server>\nprint server log",
	}
	c.commands["serverlist"] = Command{
		function: c.getServerList,
		desc:     "return name list of servers",
	}
}

func (c *Console) Run() {
	fmt.Println("Console Simulator")
	fmt.Println("------------------")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		argv := strings.Split(input, " ")

		command, ok := c.commands[argv[0]]
		if !ok {
			fmt.Printf("command not found: %s\n", argv[0])
			continue
		}
		command.function(argv)
	}
}
