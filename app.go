package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	// "go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type ComPortInfo struct {
	Name string
	Usb  string
	Vid  string
	Pid  string
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) UpdatePortList() []ComPortInfo {
	fmt.Println("back: UpdatePortList()")

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		// TODO: Вернуть ошибку на фронт.
		log.Fatal(err)
	}
	if len(ports) == 0 {
		return []ComPortInfo{}
	}

	var portList []ComPortInfo
	for _, port := range ports {
		fmt.Printf("Port: %s\n", port.Name)

		if port.Product != "" {
			fmt.Printf("   Product Name: %s\n", port.Product)
		}
		if port.IsUSB {
			fmt.Printf("   USB ID      : %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial  : %s\n", port.SerialNumber)
		}

		portList = append(portList, ComPortInfo{
			Name: port.Name,
			Usb:  strconv.FormatBool(port.IsUSB),
			Vid:  port.VID,
			Pid:  port.PID,
		})
	}

	return portList
}
