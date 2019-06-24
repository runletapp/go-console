package interfaces

import (
	"io"
)

// Console communication interface
type Console interface {
	io.Reader
	io.Writer
	io.Closer

	// SetSize sets the console size
	SetSize(cols int, rows int) error

	// GetSize gets the console size
	// cols, rows, error
	GetSize() (int, int, error)

	// Start starts the process with the supplied args
	Start(args []string) error

	// Wait waits the process to finish
	Wait() error
}
