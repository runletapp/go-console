package console

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnloadEmbeddedDeps(t *testing.T) {
	procI, err := New(120, 60)
	assert.Nil(t, err)

	proc := procI.(*consoleWindows)

	dllPath, err := proc.UnloadEmbeddedDeps()
	assert.Nil(t, err)

	files, err := os.ReadDir(dllPath)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(files))

	filenames := []string{}

	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	assert.Contains(t, filenames, "winpty-agent.exe")
	assert.Contains(t, filenames, "winpty.dll")
}
