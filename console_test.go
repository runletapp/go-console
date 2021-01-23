package console

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createSnapshot(t *testing.T, filename string, data []byte) {
	assert := assert.New(t)

	file, err := os.Create(filename)
	assert.Nil(err)
	defer file.Close()

	_, err = file.Write(data)
	assert.Nil(err)

	t.Fatalf("Snapshot created")
}

func checkSnapshot(t *testing.T, name string, data []byte) {
	assert := assert.New(t)

	snapshot := fmt.Sprintf(path.Join("snapshots", runtime.GOOS, "%s.snap"), name)
	assert.Nil(os.MkdirAll(path.Dir(snapshot), 0755))

	file, err := os.Open(snapshot)
	if err != nil {
		createSnapshot(t, snapshot, data)
	}
	defer file.Close()

	snapshotData, err := ioutil.ReadAll(file)
	assert.Nil(err)

	assert.EqualValues(snapshotData, data)
}

func TestRun(t *testing.T) {
	assert := assert.New(t)

	var args []string
	if runtime.GOOS == "windows" {
		args = []string{"echo", "windows"}
	} else {
		args = []string{"printf", "with \033[0;31mCOLOR\033[0m"}
	}

	proc, err := New(120, 60)
	assert.Nil(err)

	err = proc.Start(args)
	assert.Nil(err)
	defer proc.Close()

	data, _ := ioutil.ReadAll(proc)

	if runtime.GOOS == "windows" {
		assert.Truef(bytes.Contains(data, []byte("windows")), "Does not contain output")
	} else {
		checkSnapshot(t, "TestRun", data)
	}
}

func TestSize(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("skipped in windows")
	}

	assert := assert.New(t)

	args := []string{"stty", "size"}

	proc, err := New(120, 60)
	assert.Nil(err)

	assert.Nil(proc.Start(args))

	data, _ := ioutil.ReadAll(proc)

	assert.Truef(bytes.Contains(data, []byte("60 120")), "Does not contain size")
}

func TestSize2(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("skipped in windows")
	}

	assert := assert.New(t)

	args := []string{"stty", "size"}

	proc, err := New(60, 120)
	assert.Nil(err)

	assert.Nil(proc.Start(args))

	data, _ := ioutil.ReadAll(proc)

	assert.Truef(bytes.Contains(data, []byte("120 60")), "Does not contain size")
}

func TestWait(t *testing.T) {
	assert := assert.New(t)

	var args []string
	if runtime.GOOS == "windows" {
		args = []string{"sleep", "5s"}
	} else {
		args = []string{"sleep", "5s"}
	}

	proc, err := New(120, 60)
	assert.Nil(err)

	assert.Nil(proc.Start(args))
	defer proc.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, err := proc.Wait()
		assert.Nil(err)
		wg.Done()
	}()

	n, _ := io.Copy(os.Stdout, proc)

	var res int64
	if runtime.GOOS == "windows" {
		res = 8
	}

	assert.Equal(int64(res), n)

	wg.Wait()
}

func TestCWD(t *testing.T) {
	assert := assert.New(t)

	args := []string{"pwd"}
	if runtime.GOOS == "windows" {
		args = []string{"cmd", "/c", "echo", "%cd%"}
	}

	proc, err := New(120, 60)
	assert.Nil(err)
	defer proc.Close()

	tmpdir, err := ioutil.TempDir("", "go-console_")
	assert.Nil(err)
	defer os.RemoveAll(tmpdir)

	assert.Nil(proc.SetCWD(tmpdir))

	assert.Nil(proc.Start(args))

	data, _ := ioutil.ReadAll(proc)

	assert.Contains(string(data), tmpdir)
}

func TestENV(t *testing.T) {
	assert := assert.New(t)

	args := []string{"env"}
	if runtime.GOOS == "windows" {
		args = []string{"cmd", "/c", "echo", "MYENV=%MYENV%"}
	}

	proc, err := New(120, 60)
	assert.Nil(err)
	defer proc.Close()

	assert.Nil(proc.SetENV([]string{"MYENV=test"}))

	assert.Nil(proc.Start(args))

	data, _ := ioutil.ReadAll(proc)

	assert.Contains(string(data), "MYENV=test")
}

func TestPID(t *testing.T) {
	assert := assert.New(t)

	args := []string{"sleep", "5s"}

	proc, err := New(120, 60)
	assert.Nil(err)

	assert.Nil(proc.Start(args))
	defer proc.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, err := proc.Wait()
		assert.Nil(err)
		wg.Done()
	}()

	pid, err := proc.Pid()
	assert.Nil(err)
	assert.NotEqual(0, pid)

	n, _ := io.Copy(os.Stdout, proc)

	var res int64
	if runtime.GOOS == "windows" {
		res = 8
	}

	assert.Equal(int64(res), n)

	wg.Wait()
}

func TestKill(t *testing.T) {
	assert := assert.New(t)

	args := []string{"sleep", "1h"}

	proc, err := New(120, 60)
	assert.Nil(err)

	assert.Nil(proc.Start(args))
	defer proc.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		state, err := proc.Wait()
		assert.Nil(err)

		xSignal := "killed"
		if runtime.GOOS == "windows" {
			xSignal = "signal -1"
		}

		signal := state.Sys().(syscall.WaitStatus).Signal()
		assert.Equal(xSignal, signal.String())
		wg.Done()
	}()

	time.Sleep(1 * time.Second)
	assert.Nil(proc.Kill())

	n, _ := io.Copy(os.Stdout, proc)

	var res int64
	if runtime.GOOS == "windows" {
		res = 8
	}

	assert.Equal(int64(res), n)

	wg.Wait()
}

func TestSignal(t *testing.T) {
	assert := assert.New(t)

	args := []string{"sleep", "1h"}

	proc, err := New(120, 60)
	assert.Nil(err)

	assert.Nil(proc.Start(args))
	defer proc.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		state, err := proc.Wait()
		assert.Nil(err)

		xSignal := "killed"
		if runtime.GOOS == "windows" {
			xSignal = "signal -1"
		}

		signal := state.Sys().(syscall.WaitStatus).Signal()
		assert.Equal(xSignal, signal.String())
		wg.Done()
	}()

	time.Sleep(1 * time.Second)
	assert.Nil(proc.Signal(os.Kill))

	n, _ := io.Copy(os.Stdout, proc)

	var res int64
	if runtime.GOOS == "windows" {
		res = 8
	}

	assert.Equal(int64(res), n)

	wg.Wait()
}
