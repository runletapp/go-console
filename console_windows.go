package console

import (
	"os"
	"strings"
	"syscall"

	"github.com/runletapp/go-console/interfaces"
	"github.com/runletapp/go-winpty"
)

var _ interfaces.Console = (*consoleWindows)(nil)

type consoleWindows struct {
	initialCols int
	initialRows int

	file *winpty.WinPTY

	cwd string
	env []string
}

func newNative(cols int, rows int) (Console, error) {
	return &consoleWindows{
		initialCols: cols,
		initialRows: rows,

		file: nil,

		cwd: ".",
		env: os.Environ(),
	}, nil
}

func (c *consoleWindows) Start(args []string) error {
	opts := winpty.Options{
		InitialCols: uint32(c.initialCols),
		InitialRows: uint32(c.initialRows),
		Command:     strings.Join(args, " "),
		Dir:         c.cwd,
		Env:         c.env,
	}

	cmd, err := winpty.OpenWithOptions(opts)
	if err != nil {
		return err
	}

	c.file = cmd
	return nil
}

func (c *consoleWindows) Read(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	n, err := c.file.StdOut.Read(b)

	return n, err
}

func (c *consoleWindows) Write(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.file.StdIn.Write(b)
}

func (c *consoleWindows) Close() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	c.file.Close()
	return nil
}

func (c *consoleWindows) SetSize(cols int, rows int) error {
	c.initialRows = rows
	c.initialCols = cols

	if c.file == nil {
		return nil
	}

	c.file.SetSize(uint32(c.initialCols), uint32(c.initialRows))
	return nil
}

func (c *consoleWindows) GetSize() (int, int, error) {
	return c.initialCols, c.initialRows, nil
}

func (c *consoleWindows) Wait() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	handle := c.file.GetProcHandle()

	_, err := syscall.WaitForSingleObject(syscall.Handle(handle), syscall.INFINITE)

	return err
}

func (c *consoleWindows) SetCWD(cwd string) error {
	c.cwd = cwd
	return nil
}

func (c *consoleWindows) SetENV(environ []string) error {
	c.env = append(os.Environ(), environ...)
	return nil
}
