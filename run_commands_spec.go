package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"
)

/* ====================================================== *
 * FAKE VALUES
 * ====================================================== */

type SystemCommand struct {
	executable string
	args       []string
	stdout     string
	err        error
}

func NewSystemCommand(executable string, args ...string) Command {
	return SystemCommand{
		executable: executable,
		args:       args,
	}
}

func (s SystemCommand) Data() string {
	return s.stdout
}

func (s SystemCommand) Error() error {
	return s.err
}

func (s SystemCommand) Execute() Command {
	stdout, err := exec.Command(s.executable, s.args...).Output()
	s.stdout = string(stdout)
	s.err = err
	return s
}

/* ====================================================== *
 * TESTS
 * ====================================================== */

var _ = Describe("Running a sequence of commands", func() {
	It("does nothing given a no-op command generator", func(done Done) {
		result, err := RunCommandsFrom(func(ch chan Command) {
			close(ch)
		})
		Expect(result).To(Equal(""))
		Expect(err).To(BeNil())

		close(done)
	})

	It("runs a successful command and returns the result", func(done Done) {
		result, err := RunCommandsFrom(func(ch chan Command) {
			ch <- NewSystemCommand("echo", "hi")
			<-ch
			close(ch)
		})
		Expect(result).To(Equal("hi\n"))
		Expect(err).To(BeNil())

		close(done)
	})

	It("returns the result from the last command run", func(done Done) {
		result, err := RunCommandsFrom(func(ch chan Command) {
			ch <- NewSystemCommand("echo", "hello")
			<-ch
			ch <- NewSystemCommand("echo", "goodbye")
			<-ch
			close(ch)
		})
		Expect(result).To(Equal("goodbye\n"))
		Expect(err).To(BeNil())

		close(done)
	})

	It("returns the error from a failing command", func(done Done) {
		result, err := RunCommandsFrom(func(ch chan Command) {
			ch <- NewSystemCommand("false")
			<-ch
			close(ch)
		})
		Expect(result).To(Equal(""))
		Expect(err).To(HaveOccurred())

		close(done)
	})

	It("passes the output of a command back to the command generator", func(done Done) {
		result, err := RunCommandsFrom(func(ch chan Command) {
			ch <- NewSystemCommand("echo", "-n", "hi")
			result := <-ch
			ch <- NewSystemCommand("echo", "-n", "got result: "+result.Data())
			<-ch
			close(ch)
		})
		Expect(result).To(Equal("got result: hi"))
		Expect(err).To(BeNil())

		close(done)
	})

	It("passes the error from a command back to the command generator", func(done Done) {
		result, err := RunCommandsFrom(func(ch chan Command) {
			ch <- NewSystemCommand("false")
			result := <-ch
			if result.Error() != nil {
				ch <- NewSystemCommand("echo", "-n", "this should work")
				<-ch
			}
			close(ch)
		})
		Expect(result).To(Equal("this should work"))
		Expect(err).To(BeNil())

		close(done)
	})
})
