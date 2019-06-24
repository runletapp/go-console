# go-console

[![Build Status](https://travis-ci.com/runletapp/go-console.svg?branch=master)](https://travis-ci.com/runletapp/go-console)
[![GoDoc](https://godoc.org/github.com/runletapp/go-console?status.svg)](https://godoc.org/github.com/runletapp/go-console)

`go-console` is a cross-platform `PTY` interface. On *nix platforms we rely on [pty](https://github.com/kr/pty) and on windows [go-winpty](https://github.com/iamacarpet/go-winpty) (this means that you also need [winpty](https://github.com/rprichard/winpty))

## Example

```go
package main

import (
	"io"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/runletapp/go-console"
)

func main() {

	proc, err := console.New(120, 60)
	if err != nil {
		panic(err)
	}
	defer proc.Close()

	var args []string

	if runtime.GOOS == "windows" {
		args = []string{"cmd.exe", "/c", "dir"}
	} else {
		args = []string{"ls", "-lah", "--color"}
	}

	if err := proc.Start(args); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		_, err = io.Copy(os.Stdout, proc)
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
	}()

	if err := proc.Wait(); err != nil {
		log.Printf("Wait err: %v\n", err)
	}

	wg.Wait()
}

```