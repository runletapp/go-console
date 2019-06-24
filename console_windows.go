package console

import (
	"log"
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
}

func newNative(cols int, rows int) (Console, error) {
	return &consoleWindows{
		initialCols: cols,
		initialRows: rows,

		file: nil,
	}, nil
}

func (c *consoleWindows) Start(args []string) error {
	opts := winpty.Options{
		InitialCols: uint32(c.initialCols),
		InitialRows: uint32(c.initialRows),
		Command:     strings.Join(args, " "),
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
	log.Printf("Handle: %v", handle)

	_, err := syscall.WaitForSingleObject(syscall.Handle(handle), syscall.INFINITE)

	return err
}
