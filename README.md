# Mock-Free Testing of Side Effects

This is an example repo that shows you how to test code that
has side effects (i.e. depends on system calls)...

- without actually making those system calls
- without writing unreadable, mock-heavy tests

The made-up problem the example is solving is this: say your
app needs to read config from a file. If the file doesn't
exist, it should fall back to a default config file
location. If the default file doesn't exist either, it
should error out.

The question is: *how do you test this logic without
creating any files or using any test doubles*? Read the code
and find out!

## Setup

You'll need `ginkgo` and `gomega`

```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega/...
```

Ensure that `$GOPATH/bin` is on your `$PATH`. If it is,
`which ginkgo` should show you the path to the `ginkgo`
executable.

## Running the tests

```
ginkgo -r
```

## Running the program

```
go run main.go
```

If you want it to output anything other than "couldn't read
config file", do one or both of the following:

```
echo custom config > /tmp/.my-app.cfg
```

```
echo default config > /tmp/.my-app.default.cfg
```

## A note on design

This example is obviously way over-engineered. This repo
spends 100+ lines of code and 200+ lines of tests to *read a
config file*. That's because it's just an example. You
wouldn't want to use this technique in real code unless you
were doing something more complex.
