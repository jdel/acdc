package lgr

import (
	"fmt"
	"io"

	"github.com/docker/libcompose/logger"
)

// ACDCLogger is a logger.Logger and logger.Factory implementation that prints raw data with no formatting.
type ACDCLogger struct {
	Output chan []byte
}

// Out is a no-op function.
func (r *ACDCLogger) Out(message []byte) {
	// fmt.Print(string(message))
	output := string(message)
	fmt.Print(output)
	go func() {
		r.Output <- message
		close(r.Output)
	}()
}

// Err is a no-op function.
func (r *ACDCLogger) Err(message []byte) {
	// fmt.Fprint(os.Stderr, string(message))
	go func() {
		r.Output <- message
		close(r.Output)
	}()
}

// CreateContainerLogger allows ACDCLogger to implement logger.Factory.
func (r *ACDCLogger) CreateContainerLogger(_ string) logger.Logger {
	return &ACDCLogger{}
}

// CreateBuildLogger allows ACDCLogger to implement logger.Factory.
func (r *ACDCLogger) CreateBuildLogger(_ string) logger.Logger {
	return &ACDCLogger{}
}

// CreatePullLogger allows ACDCLogger to implement logger.Factory.
func (r *ACDCLogger) CreatePullLogger(_ string) logger.Logger {
	return &ACDCLogger{}
}

// OutWriter returns the base writer
func (r *ACDCLogger) OutWriter() io.Writer {
	// return os.Stdout
	return nil
}

// ErrWriter returns the base writer
func (r *ACDCLogger) ErrWriter() io.Writer {
	// return os.Stderr
	return nil
}
