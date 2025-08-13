package main

import (
	"fmt"
	"log"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func main() {
	port_names, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(port_names) == 0 {
		log.Fatal("No serial ports found!")
	}
	for i, port_name := range port_names {
		fmt.Printf("Port[%v]: %v\n", i, port_name)
	}

	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(port_names[1], mode)
	if err != nil {
		log.Fatal(err)
	}
	_ = port

	/*
		n, err := port.Write([]byte("10,20,30\n\r"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %v bytes\n", n)

		buff := make([]byte, 100)
		for {
			n, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
				break
			}
			if n == 0 {
				fmt.Println("\nEOF")
				break
			}
			fmt.Printf("%v", string(buff[:n]))
		}
	*/

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return
	}
	for _, port := range ports {
		fmt.Printf("Found port: %s\n", port.Name)
		if port.IsUSB {
			fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial %s\n", port.SerialNumber)
		}
	}
}
