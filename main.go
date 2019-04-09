package main

import (
	"fmt"
	"io/ioutil"
)

/* ====================================================== *
 * ENTRYPOINT
 * ====================================================== */

func main() {
	config, err := GetConfig()
	if err != nil {
		fmt.Println("couldn't read config file")
	} else {
		fmt.Println("config file contents: " + config)
	}
}

/* ====================================================== *
 * TYPE DECLARATIONS
 * ====================================================== */

type Command interface {
	Data() string
	Error() error
	Execute() Command
}

type NullCommand struct {
}

func (NullCommand) Data() string {
	return ""
}

func (NullCommand) Error() error {
	return nil
}

func (n NullCommand) Execute() Command {
	return n
}

type ReadFileCommand struct {
	Path         string
	fileContents string
	err          error
}

func (r ReadFileCommand) Data() string {
	return r.fileContents
}

func (r ReadFileCommand) Error() error {
	return r.err
}

func (r ReadFileCommand) Execute() Command {
	bytes, err := ioutil.ReadFile(r.Path)
	r.err = err
	r.fileContents = string(bytes)
	return r
}

/* ====================================================== *
 * ROUTINES
 * ====================================================== */

func GetConfig() (string, error) {
	return RunCommandsFrom(ReadConfigFiles)
}

func RunCommandsFrom(generator func(Await)) (string, error) {
	ch := StartCoroutine(generator)

	var result Command = NullCommand{}
	for cmd := range ch {
		result = cmd.Execute()
		ch <- result
	}
	return result.Data(), result.Error()
}

func StartCoroutine(generator func(Await)) chan Command {
	ch := make(chan Command)
	go func() {
		defer close(ch)
		generator(makeAwaitFunc(ch))
	}()
	return ch
}

func ReadConfigFiles(await Await) {
	result := await(ReadFileCommand{Path: "/tmp/.my-app.cfg"})
	if result.Error() == nil {
		return
	}

	await(ReadFileCommand{Path: "/tmp/.my-app.default.cfg"})
}

type Await func(Command) Command

func makeAwaitFunc(ch chan Command) Await {
	return func(cmd Command) Command {
		ch <- cmd
		return <-ch
	}
}
