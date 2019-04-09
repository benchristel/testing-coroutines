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

var _ = Describe("Running a coroutine", func() {
	It("does nothing given a no-op command generator", func(done Done) {
		result, err := RunCoroutine(func(_ Await) {})
		Expect(result).To(Equal(""))
		Expect(err).To(BeNil())

		close(done)
	})

	It("runs a successful command and returns the result", func(done Done) {
		result, err := RunCoroutine(func(await Await) {
			await(NewSystemCommand("echo", "hi"))
		})
		Expect(result).To(Equal("hi\n"))
		Expect(err).To(BeNil())

		close(done)
	})

	It("returns the result from the last command run", func(done Done) {
		result, err := RunCoroutine(func(await Await) {
			await(NewSystemCommand("echo", "hello"))
			await(NewSystemCommand("echo", "goodbye"))
		})
		Expect(result).To(Equal("goodbye\n"))
		Expect(err).To(BeNil())

		close(done)
	})

	It("returns the error from a failing command", func(done Done) {
		result, err := RunCoroutine(func(await Await) {
			await(NewSystemCommand("false"))
		})
		Expect(result).To(Equal(""))
		Expect(err).To(HaveOccurred())

		close(done)
	})

	It("passes the output of a command back to the command generator", func(done Done) {
		result, err := RunCoroutine(func(await Await) {
			result := await(NewSystemCommand("echo", "-n", "hi"))
			await(NewSystemCommand("echo", "-n", "got result: "+result.Data()))
		})
		Expect(result).To(Equal("got result: hi"))
		Expect(err).To(BeNil())

		close(done)
	})

	It("passes the error from a command back to the command generator", func(done Done) {
		result, err := RunCoroutine(func(await Await) {
			result := await(NewSystemCommand("false"))
			if result.Error() != nil {
				await(NewSystemCommand("echo", "-n", "this should work"))
			}
		})
		Expect(result).To(Equal("this should work"))
		Expect(err).To(BeNil())

		close(done)
	})
})
