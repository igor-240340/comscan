package main

import (
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

// Реализует интерфейс ComPort.
type ComPortReal struct {
	port serial.Port
}

func (r *ComPortReal) Enumerate() ([]*enumerator.PortDetails, error) {
	return enumerator.GetDetailedPortsList()
}

func (r *ComPortReal) Open(portName string) error {
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(portName, mode)
	r.port = port
	return err
}

func (r *ComPortReal) Close() error {
	return r.port.Close()
}

func (r *ComPortReal) Read(buf []byte) (int, error) {
	n, err := r.port.Read(buf)
	return n, err
}

func (r *ComPortReal) Write(b []byte) (int, error) {
	return r.port.Write(b)
}
