package main

import (
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

// Реализует интерфейс ComPort и
// использует для реализации библиотеку serial.
//
// NOTE: У текущей реализации, возможно, есть недостаток.
// Поскольку port спрятан внутри структуры, мы сейчас никак не защищаем код
// от, скажем, двойного вызова serial.Open, что приведет к переопределению port.
// Таким образом можем потерять "указатель" на уже открытый порт и не сможем его корректно закрыть.
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
