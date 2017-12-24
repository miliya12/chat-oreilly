package trace

import (
	"fmt"
	"io"
)

// Tracer is interface express object available on recording event in codes
type Tracer interface {
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

// disable logger

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}
func Off() Tracer {
	return &nilTracer{}
}
