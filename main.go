package main

import (
	"fmt"
	"io/ioutil"
)

type Command interface {
	Data() string
	Error() error
	Execute() Command
}

type ReadFileCommand struct {
	Path         string
	fileContents string
	err          error
}

func (r ReadFileCommand) Execute() Command {
	bytes, err := ioutil.ReadFile(r.Path)
	r.err = err
	r.fileContents = string(bytes)
	return r
}

func (r ReadFileCommand) Error() error {
	return r.err
}

func (r ReadFileCommand) Data() string {
	return r.fileContents
}

func main() {
	config, err := GetConfig()
	if err != nil {
		fmt.Println("couldn't read config file")
	} else {
		fmt.Println("Config file contents: " + config)
	}
}

func GetConfig() (string, error) {
	return RunCmdsFrom(ReadConfigFiles)
}

func RunCmdsFrom(generator func(chan Command)) (string, error) {
	ch := make(chan Command)
	go generator(ch)

	var cmd, result Command
	for cmd = range ch {
		result = cmd.Execute()
		ch <- result
	}
	return result.Data(), result.Error()
}

func ReadConfigFiles(ch chan Command) {
	defer close(ch)

	result := await(ch, ReadFileCommand{Path: "/tmp/.my-app.cfg"})
	if result.Error() == nil {
		return
	}

	await(ch, ReadFileCommand{Path: "/tmp/.my-app.default.cfg"})
}

func await(ch chan Command, cmd Command) Command {
	ch <- cmd
	return <-ch
}
