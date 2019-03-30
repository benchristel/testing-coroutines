# Mock-Free Testing of Coroutines

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
