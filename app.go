package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// Информация о порте, которая уходит на фронт.
type ComPortInfo struct {
	Name         string
	Usb          string
	Vid          string
	Pid          string
	SentData     string // Строки ping/pong,
	ReceivedData string // для портов, подпадающих по условие из ТЗ.
}

type App struct {
	ctx     context.Context
	comport ComPort
}

func NewApp(comport ComPort) *App {
	return &App{comport: comport}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Получает новый список портов.
// А для тех, которые подпадают под условие из ТЗ, делает пинг.
//
// NOTE: Здесь для простоты пишем/читаем порты синхронно/последовательно.
// Но можно попробовать открыть и запустить запись/чтение сразу по N портам параллельно.
// Затем ждать результатов по всем портам по аналогии с 'await Promise.all()'.
// Я знаю, что в Go есть горутины, но я пока с ними не разбирался.
// Хочу покрыть этот метод юнит-тестами, поэтому навряд ли доберусь до горутин в рамках тестового.
// Для начала пускай будет не самая лучшая, но хотя бы в какой-то степени корректная реализация.
// Да и есть ли смысл параллелить, если одновременно будет около десятка устройств,
// здесь надо смотреть, какой объем данных нужно читать/писать на устройства.
func (a *App) UpdatePortList() ([]ComPortInfo, error) {
	const sendMessage string = "AT+VERSION\r\n"

	portInfos, err := a.comport.Enumerate()
	if err != nil {
		return nil, err
	}
	if len(portInfos) == 0 {
		return []ComPortInfo{}, err
	}

	var portList []ComPortInfo
	for i, portInfo := range portInfos {
		portList = append(portList, ComPortInfo{
			Name:         portInfo.Name,
			Usb:          strconv.FormatBool(portInfo.IsUSB),
			Vid:          portInfo.VID,
			Pid:          portInfo.PID,
			SentData:     "",
			ReceivedData: "",
		})

		vid := strings.ToLower(portInfo.VID)
		pid := strings.ToLower(portInfo.PID)
		if vid == "2e8a" && (pid == "f00a" || pid == "f00f") {
			fmt.Printf("Open port: %s\n", portInfo.Name)

			err := a.comport.Open(portInfo.Name)
			if err != nil {
				return nil, err
			}

			n, err := a.comport.Write([]byte(sendMessage))
			if err != nil {
				return nil, err
			}
			portList[i].SentData = sendMessage
			fmt.Printf("Sent %v bytes\n", n)

			buff := make([]byte, 100)
			var msg []byte
			for {
				n, err := a.comport.Read(buff)
				if err != nil {
					return nil, err
				}
				fmt.Printf("%s", string(buff[:n]))

				msg = append(msg, buff[:n]...)

				if strings.Contains(string(msg), "\r\n") {
					portList[i].ReceivedData = strings.SplitAfterN(string(msg), "\r\n", 2)[0]
					break
				}
			}
			err = a.comport.Close()
			if err != nil {
				return nil, err
			}
		}
	}

	return portList, err
}
