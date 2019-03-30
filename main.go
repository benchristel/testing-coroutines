package main

import (
	"fmt"
	"io/ioutil"
)

type Command interface {
	Execute() Command
	Error() error
	Result() string
}

type ReadFileCommand struct {
	Path         string
	fileContents string
	err          error
}

func (r *ReadFileCommand) Execute() Command {
	bytes, err := ioutil.ReadFile(r.Path)
	r.err = err
	r.fileContents = string(bytes)
	return r
}

func (r *ReadFileCommand) Error() error {
	return r.err
}

func (r *ReadFileCommand) Result() string {
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
	var cmd Command
	go generator(ch)
	for cmd = range ch {
		ch <- cmd.Execute()
	}
	return cmd.Result(), cmd.Error()
}

func ReadConfigFiles(ch chan Command) {
	defer close(ch)

	ch <- &ReadFileCommand{Path: "/tmp/.my-app.cfg"}
	result := <-ch
	if result.Error() == nil {
		return
	}

	ch <- &ReadFileCommand{Path: "/tmp/.my-app.default.cfg"}
	_ = <-ch
}
