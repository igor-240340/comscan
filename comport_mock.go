package main

import (
	"errors"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type ComPortMock struct {
	port serial.Port
}

// Реализует интерфейс ComPort.
func (r *ComPortMock) Enumerate() ([]*enumerator.PortDetails, error) {
	return nil, errors.New("ERROR: ComPortMock.Enumerate()")
}

func (r *ComPortMock) Open(portName string) error {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(portName, mode)
	r.port = port
	return err
}

func (r *ComPortMock) Close() error {
	return r.port.Close()
}

func (r *ComPortMock) Read(buf []byte) (int, error) {
	n, err := r.port.Read(buf)
	return n, err
}

func (r *ComPortMock) Write(b []byte) (int, error) {
	return r.port.Write(b)
}
