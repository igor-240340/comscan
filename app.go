package main

import (
	"context"
	"fmt"
	"log"
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

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Получает новый список портов.
// А для тех, которые подпадают под условие из ТЗ, делает пинг.
//
// NOTE: Здесь для простоты пишем/читаем порты последовательно, но
// можно попробовать открыть и запустить запись/чтение сразу по N портам параллельно,
// дожидаясь результатов по всем портам по аналогии с 'await Promise.all()'.
// Я знаю, что в Go есть корутины, но я пока с ними не разбирался.
// Хочу написать юнит-тесты, поэтому навряд ли доберусь до корутин.
// Для начала пускай будет не самая лучшая, но хотя бы в какой-то степени корректная реализация.
func (a *App) UpdatePortList() ([]ComPortInfo, error) {
	fmt.Println("back: UpdatePortList()")

	portInfos, err := a.comport.Enumerate()
	if err != nil {
		// log.Fatal(err)
		return []ComPortInfo{}, fmt.Errorf("something went wrong")
	}
	if len(portInfos) == 0 {
		return []ComPortInfo{}, err
	}

	var portList []ComPortInfo
	// var selectedPortList []ComPortInfo
	for i, portInfo := range portInfos {
		fmt.Printf("Port: %s\n", portInfo.Name)

		if portInfo.Product != "" {
			fmt.Printf("   Product Name: %s\n", portInfo.Product)
		}
		if portInfo.IsUSB {
			fmt.Printf("   USB ID      : %s:%s\n", portInfo.VID, portInfo.PID)
			fmt.Printf("   USB serial  : %s\n", portInfo.SerialNumber)
		}

		portList = append(portList, ComPortInfo{
			Name:         portInfo.Name,
			Usb:          strconv.FormatBool(portInfo.IsUSB),
			Vid:          portInfo.VID,
			Pid:          portInfo.PID,
			SentData:     "",
			ReceivedData: "",
		})

		// Условие.
		var strBuff string
		// const VID uint16 = 0x2e8a
		// const PID1 uint16 = 0xf00a
		// const PID2 uint16 = 0xf00f
		const VID uint16 = 0x0193
		const PID1 uint16 = 0x1771
		const PID2 uint16 = 0xf00f
		// if (port.VID == strconv.Itoa(int(VID))) && (port.PID == strconv.Itoa(int(PID1)) || port.PID == strconv.Itoa(int(PID2))) {
		if portInfo.Name == "COM2" {
			err := a.comport.Open(portInfo.Name)
			if err != nil {
				log.Fatal(err)
			}

			n, err := a.comport.Write([]byte("AT+VERSION\r\n"))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Sent %v bytes\n", n)

			buff := make([]byte, 100)
			for {
				// Reads up to 100 bytes
				n, err := a.comport.Read(buff)
				if err != nil {
					log.Fatal(err)
				}
				if n == 0 {
					fmt.Println("\nEOF")
					break
				}

				fmt.Printf("%s", string(buff[:n]))
				strBuff = string(buff[:n])

				// If we receive a newline stop reading
				if strings.Contains(string(buff[:n]), "\n") {
					break
				}
			}
			a.comport.Close()

			portList[i].SentData = "AT+VERSION\r\n"
			portList[i].ReceivedData = strBuff
		}
	}

	return portList, err
}
