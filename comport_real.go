package main

import (
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

// Реализует интерфейс ComPort и
// использует для реализации библиотеку serial.
//
// NOTE: У этой реализации, возможно, есть недостаток.
// Поскольку port спрятан внутри структуры, то код сейчас никак не защищён
// от, скажем, двойного вызова serial.Open, что приведет к перезаписи port и
// "потере" указателя на уже открытый порт, а это значит, что мы не сможем его корректно закрыть.
//
// Но в данном приложении такой сценарий не возможен, поскольку на стороне фронта
// мы обеспечиваем синхронное взаимодействие с бэком и параллельных вызовов бэка не будет.
type ComPortReal struct {
	port serial.Port
}

func (с *ComPortReal) Enumerate() ([]*enumerator.PortDetails, error) {
	return enumerator.GetDetailedPortsList()
}

func (с *ComPortReal) Open(portName string) error {
	// NOTE: Конечно, для чистоты, конфигуратор подключения нужно вынести наружу
	// и передавать нашу собственную структуру, чтобы не привязываться к конкретной библиотеке.
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(portName, mode)
	с.port = port
	return err
}

func (с *ComPortReal) Close() error {
	return с.port.Close()
}

func (с *ComPortReal) Read(buf []byte) (int, error) {
	n, err := с.port.Read(buf)
	return n, err
}

func (с *ComPortReal) Write(b []byte) (int, error) {
	return с.port.Write(b)
}
