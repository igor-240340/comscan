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
	SentData     string // Строки ping/pong
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
	fmt.Println("back: UpdatePortList()")

	portInfos, err := a.comport.Enumerate()
	if err != nil {
		return nil, err
	}
	if len(portInfos) == 0 {
		return []ComPortInfo{}, err
	}

	var portList []ComPortInfo
	for i, portInfo := range portInfos {
		/*
			fmt.Printf("Port: %s\n", portInfo.Name)

			if portInfo.Product != "" {
				fmt.Printf("   Product Name: %s\n", portInfo.Product)
			}
			if portInfo.IsUSB {
				fmt.Printf("   USB ID      : %s:%s\n", portInfo.VID, portInfo.PID)
				fmt.Printf("   USB serial  : %s\n", portInfo.SerialNumber)
			}
		*/

		portList = append(portList, ComPortInfo{
			Name:         portInfo.Name,
			Usb:          strconv.FormatBool(portInfo.IsUSB),
			Vid:          portInfo.VID,
			Pid:          portInfo.PID,
			SentData:     "",
			ReceivedData: "",
		})

		// Условие.
		// const VID uint16 = 0x0193
		// const PID1 uint16 = 0x1771
		// const PID2 uint16 = 0xf00f
		// var strBuff string
		const vid uint16 = 0x2e8a  // Prod
		const pid1 uint16 = 0xf00a // Prod
		const pid2 uint16 = 0xf00f // Prod
		cond := portInfo.VID == strconv.Itoa(int(vid)) &&
			(portInfo.PID == strconv.Itoa(int(pid1)) || portInfo.PID == strconv.Itoa(int(pid2)))
			// if portInfo.Name == "COM255" {
		if cond {
			err := a.comport.Open(portInfo.Name)
			// err := a.comport.Open("sdfdsf")
			if err != nil {
				return nil, err
			}

			n, err := a.comport.Write([]byte("AT+VERSION\r\n"))
			if err != nil {
				return nil, err
			}
			fmt.Printf("Sent %v bytes\n", n)

			buff := make([]byte, 100)
			var msg []byte
			for {
				n, err := a.comport.Read(buff)
				if err != nil {
					return nil, err
				}

				msg = append(msg, buff[:n]...)
				fmt.Printf("%s", string(buff[:n]))

				if strings.Contains(string(msg), "\r\n") {
					break
				}
			}
			err = a.comport.Close()
			if err != nil {
				return nil, err
			}

			portList[i].SentData = "AT+VERSION\r\n"
			portList[i].ReceivedData = string(msg)
		}
	}

	return portList, err
}
