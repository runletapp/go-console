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
