package timoni

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Timoni struct {
	tasks map[string]Task
}

type Task struct {
	Process Process
	Server  Server
}

type Process struct {
	Stdout  io.Writer
	Stdin   io.WriteCloser
	Restart bool
}

type Server struct {
	Name    string
	Command string
	Args    string
	Dir     string
}

func New() *Timoni {
	return &Timoni{
		tasks: make(map[string]Task),
	}
}

func (t *Timoni) AddServer(server Server) {
	t.tasks[server.Name] = Task{
		Server: server,
	}
}

func (t *Timoni) SetAutorestart(name string, status bool) error {
	task, err := t.GetTask(name)
	if err != nil {
		return err
	}
	task.Process.Restart = status
	return nil
}

func (t *Timoni) GetLog(name string) (string, error) {
	task, err := t.GetTask(name)
	if err != nil {
		return "", err
	}
	if task.Process.Stdout == nil {
		return "", fmt.Errorf("server is not running")
	}
	return fmt.Sprint(task.Process.Stdout), nil
}

func (t *Timoni) SendCommand(name, command string) error {
	task, err := t.GetTask(name)
	if err != nil {
		return err
	}
	if task.Process.Stdin == nil {
		return fmt.Errorf("server is not running")
	}
	_, err = io.WriteString(task.Process.Stdin, command+"\n")
	return err
}

func (t *Timoni) GetTask(name string) (Task, error) {
	if task, ok := t.tasks[name]; ok {
		return task, nil
	}
	return Task{}, fmt.Errorf("Server %s not found", name)
}

func (t *Timoni) GetTaskList() []string {
	var list []string
	for key := range t.tasks {
		list = append(list, key)
	}
	return list
}

func (t *Timoni) Run(name string) error {
	task, err := t.GetTask(name)
	if err != nil {
		return err
	}
	cmd := exec.Command(task.Server.Command)
	cmd.Dir = task.Server.Dir
	cmd.Args = strings.Split(task.Server.Args, " ")

	var outbuf, errbuf strings.Builder
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	t.tasks[task.Server.Name] = Task{
		Server: task.Server,
		Process: Process{
			Stdout:  cmd.Stdout,
			Stdin:   stdin,
			Restart: true,
		},
	}
	task, err = t.GetTask(task.Server.Name)
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("Server %s dead with error %s. Restart!\n", task.Server.Name, err)
			timestamp := time.Now().Unix()
			err := os.WriteFile(fmt.Sprintf("./%d.%s.log", timestamp, task.Server.Name), []byte(outbuf.String()), 0644)
			if err != nil {
				println(err)
			}
			if task.Process.Restart {
				t.Run(task.Server.Name)
			}
		}
	}
	return nil
}
