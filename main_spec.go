package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
)

/* ====================================================== *
 * FAKE VALUES
 * ====================================================== */

type FailedCommand struct {
}

func (FailedCommand) Data() string {
	return ""
}

func (FailedCommand) Error() error {
	return errors.New("command failed")
}

func (f FailedCommand) Execute() Command {
	return f
}

type SuccessfulCommand struct {
}

func (s SuccessfulCommand) Data() string {
	return ""
}

func (SuccessfulCommand) Error() error {
	return nil
}

func (s SuccessfulCommand) Execute() Command {
	return s
}

/* ====================================================== *
 * TESTS
 * ====================================================== */

var _ = Describe("Reading config files", func() {
	It("errors when neither config file is readable", func(done Done) {
		ch := make(chan Command)
		go ReadConfigFiles(ch)

		Expect(<-ch).To(Equal(ReadFileCommand{
			Path: "/tmp/.my-app.cfg",
		}))

		ch <- FailedCommand{}

		Expect(<-ch).To(Equal(ReadFileCommand{
			Path: "/tmp/.my-app.default.cfg",
		}))

		ch <- FailedCommand{}

		Expect(<-ch).To(BeNil())
		Expect(ch).To(BeClosed())

		close(done)
	})

	It("reads only the custom config if it is present", func(done Done) {
		ch := make(chan Command)
		go ReadConfigFiles(ch)

		Expect(<-ch).To(Equal(ReadFileCommand{
			Path: "/tmp/.my-app.cfg",
		}))

		ch <- SuccessfulCommand{}

		Expect(<-ch).To(BeNil())
		Expect(ch).To(BeClosed())

		close(done)
	})

	It("falls back to the default config", func(done Done) {
		ch := make(chan Command)
		go ReadConfigFiles(ch)

		Expect(<-ch).To(Equal(ReadFileCommand{
			Path: "/tmp/.my-app.cfg",
		}))

		ch <- FailedCommand{}

		Expect(<-ch).To(Equal(ReadFileCommand{
			Path: "/tmp/.my-app.default.cfg",
		}))

		ch <- SuccessfulCommand{}

		Expect(<-ch).To(BeNil())
		Expect(ch).To(BeClosed())

		close(done)
	})
})
