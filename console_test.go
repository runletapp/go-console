package console

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sync"
	"testing"

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

	snapshot := fmt.Sprintf(path.Join("snapshots", runtime.GOOS, "TestRun.snap"))
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
		args = []string{"cmd.exe", "/c", "echo windows"}
	} else {
		args = []string{"printf", "with \033[0;31mCOLOR\033[0m"}
	}

	proc, err := New(120, 60)
	assert.Nil(err)

	assert.Nil(proc.Start(args))
	defer proc.Close()

	data, _ := ioutil.ReadAll(proc)

	checkSnapshot(t, "TestRun", data)
}

func TestSize(t *testing.T) {
	assert := assert.New(t)

	args := []string{"go", "run", path.Join("testing", "size", "size.go")}

	proc, err := New(120, 60)
	assert.Nil(err)

	assert.Nil(proc.Start(args))

	data, _ := ioutil.ReadAll(proc)

	assert.Equal("120 60\r\n", string(data))
}

func TestSize2(t *testing.T) {
	assert := assert.New(t)

	args := []string{"go", "run", path.Join("testing", "size", "size.go")}

	proc, err := New(60, 120)
	assert.Nil(err)

	assert.Nil(proc.Start(args))

	data, _ := ioutil.ReadAll(proc)

	assert.Equal("60 120\r\n", string(data))
}

func TestWait(t *testing.T) {
	assert := assert.New(t)

	var args []string
	if runtime.GOOS == "windows" {
		args = []string{"cmd.exe", "/c", "timeout /t 5"}
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
		proc.Wait()
		wg.Done()
	}()

	wg.Wait()
}
