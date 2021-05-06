package connectors

import (
	"fmt"

	"github.com/microlib/simple"
)

// Connections struct - all backend connections in a common object
type Connections struct {
	log *simple.Logger
}

func (r *Connections) Error(msg string, val ...interface{}) {
	r.log.Error(fmt.Sprintf(msg, val...))
}

func (r *Connections) Info(msg string, val ...interface{}) {
	r.log.Info(fmt.Sprintf(msg, val...))
}

func (r *Connections) Debug(msg string, val ...interface{}) {
	r.log.Debug(fmt.Sprintf(msg, val...))
}

func (r *Connections) Trace(msg string, val ...interface{}) {
	r.log.Trace(fmt.Sprintf(msg, val...))
}

// NewClientConnectors returns Connectors struct
func NewClientConnections(logger *simple.Logger) Clients {
	conns := &Connections{log: logger}
	return conns
}
