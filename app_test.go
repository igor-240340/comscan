package main

import (
	"testing"
)

func TestUpdatePortList_EnumerateFails(t *testing.T) {
	comport := &ComPortMock{}
	app := NewApp(comport)

	portList, err := app.UpdatePortList()
	if err == nil {
		t.Fatalf("expected fail")
	}
	if len(portList) != 0 {
		t.Fatalf("expected fail")
	}
	// fmt.Println(err)

	/*
		portInfos, err := serial.Enumerate()
		if err != nil {
			log.Fatal(err)
		}

		for _, portInfo := range portInfos {
			fmt.Printf("Port: %s\n", portInfo.Name)

			if portInfo.Product != "" {
				fmt.Printf("   Product Name: %s\n", portInfo.Product)
			}
			if portInfo.IsUSB {
				fmt.Printf("   USB ID      : %s:%s\n", portInfo.VID, portInfo.PID)
				fmt.Printf("   USB serial  : %s\n", portInfo.SerialNumber)
			}
		}
	*/

	// actual := app.UpdatePortList()
	// fmt.Printf("actual res: %s\n", actual)
}
