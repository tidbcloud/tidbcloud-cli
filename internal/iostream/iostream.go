package iostream

import (
	"io"
	"os"
)

type IOStreams struct {
	Out io.Writer
	Err io.Writer
}

func System() *IOStreams {
	return &IOStreams{
		Out: os.Stdout,
		Err: os.Stderr,
	}
}
