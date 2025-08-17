package main

import (
	"go.bug.st/serial/enumerator"
)

type ComPort interface {
	Enumerate() ([]*enumerator.PortDetails, error)
	Open(portName string) error
	Close() error
	Write([]byte) (int, error)
	Read([]byte) (int, error)
}
