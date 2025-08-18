package main

import (
	"go.bug.st/serial/enumerator"
)

type ComPort interface {
	// NOTE: Не совсем "чисто", должен возвращать нашу структуру вместо []*enumerator.PortDetails
	// чтобы не было привязки к конкретной библиотеке.
	Enumerate() ([]*enumerator.PortDetails, error)
	Open(portName string) error
	Close() error
	Write(p []byte) (int, error)
	Read(p []byte) (int, error)
}
