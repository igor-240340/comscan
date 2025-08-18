package main

import (
	"go.bug.st/serial/enumerator"
)

// Реализует интерфейс ComPort.
type ComPortMock struct {
	// Конфигураторы поведения основных методов.
	// Настраиваются в вызывающем коде: например, в тестах.
	EnumerateFunc func() ([]*enumerator.PortDetails, error)
	OpenFunc      func(portName string) error
	CloseFunc     func() error
	ReadFunc      func(p []byte) (n int, err error)
	WriteFunc     func(p []byte) (n int, err error)
}

func (c *ComPortMock) Enumerate() ([]*enumerator.PortDetails, error) {
	return c.EnumerateFunc()
}

func (c *ComPortMock) Open(portName string) error {
	return c.OpenFunc(portName)
}

func (c *ComPortMock) Close() error {
	return c.CloseFunc()
}

func (c *ComPortMock) Read(p []byte) (int, error) {
	return c.ReadFunc(p)
}

func (c *ComPortMock) Write(p []byte) (int, error) {
	return c.WriteFunc(p)
}
