package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"go.bug.st/serial/enumerator"
)

// Ошибка при вызове Enumerate().
//
// Должен вернуть на фронт nil и ошибку.
func TestUpdatePortList_Enumerate_Fail(t *testing.T) {
	comport := &ComPortMock{}
	app := NewApp(comport)

	comport.EnumerateFunc = func() ([]*enumerator.PortDetails, error) {
		return nil, fmt.Errorf("Enumerate failed")
	}

	portList, err := app.UpdatePortList()
	if err == nil {
		t.Fatalf("expected error")
	}
	if portList != nil {
		t.Fatalf("expected nil")
	}
}

// Портов нет.
//
// Должен вернуть на фронт пустой массив ComPortInfo.
func TestUpdatePortList_Enumerate_NoPorts(t *testing.T) {
	comport := &ComPortMock{}
	app := NewApp(comport)

	comport.EnumerateFunc = func() ([]*enumerator.PortDetails, error) {
		var res []*enumerator.PortDetails
		return res, nil
	}

	portList, err := app.UpdatePortList()
	if err != nil {
		t.Fatalf("expected nil")
	}
	if len(portList) > 0 {
		t.Fatalf("expected empty")
	}
}

// Есть два обычных COM-порта.
// USB-портов нет.
//
// Должен вернуть на фронт информацию только по обычным портам (без VID/PID).
// Не должен подключаться и делать пинг.
func TestUpdatePortList_Enumerate_NoUsbPorts(t *testing.T) {
	comport := &ComPortMock{}
	app := NewApp(comport)

	comport.EnumerateFunc = func() ([]*enumerator.PortDetails, error) {
		// Все остальные поля - в нули.
		port1 := &enumerator.PortDetails{Name: "COM1"}
		port2 := &enumerator.PortDetails{Name: "COM2"}

		var res []*enumerator.PortDetails
		res = append(res, port1, port2)

		return res, nil
	}

	expectedPortList := []ComPortInfo{
		{Name: "COM1", Usb: "false", Vid: "", Pid: "", SentData: "", ReceivedData: ""},
		{Name: "COM2", Usb: "false", Vid: "", Pid: "", SentData: "", ReceivedData: ""},
	}
	actualPortList, err := app.UpdatePortList()
	if err != nil {
		t.Fatalf("expected nil")
	}
	if len(actualPortList) != len(expectedPortList) {
		t.Fatalf("expected %d", len(expectedPortList))
	}
	if !reflect.DeepEqual(actualPortList, expectedPortList) {
		t.Fatalf("not equal: \nactual=%v\nexpected=%v", actualPortList, expectedPortList)
	}
}

// Есть два обычных COM-порта.
// Есть один USB COM-порт, но его VID/PID не удовлетворяет условию.
//
// Должен вернуть на фронт информацию по всем портам (включая VID/PID для USB-порта).
// Не должен подключаться и делать пинг.
func TestUpdatePortList_UsbPortConditionFalse(t *testing.T) {
	comport := &ComPortMock{}
	app := NewApp(comport)

	comport.EnumerateFunc = func() ([]*enumerator.PortDetails, error) {
		// Все остальные поля - в нули.
		port1 := &enumerator.PortDetails{Name: "COM1"}
		port2 := &enumerator.PortDetails{Name: "COM2"}
		port3 := &enumerator.PortDetails{Name: "COM3", IsUSB: true, VID: "0403", PID: "6001"}

		var res []*enumerator.PortDetails
		res = append(res, port1, port2, port3)

		return res, nil
	}

	expectedPortList := []ComPortInfo{
		{Name: "COM1", Usb: "false", Vid: "", Pid: "", SentData: "", ReceivedData: ""},
		{Name: "COM2", Usb: "false", Vid: "", Pid: "", SentData: "", ReceivedData: ""},
		{Name: "COM3", Usb: "true", Vid: "0403", Pid: "6001", SentData: "", ReceivedData: ""},
	}
	actualPortList, err := app.UpdatePortList()
	if err != nil {
		t.Fatalf("expected nil")
	}
	if len(actualPortList) != len(expectedPortList) {
		t.Fatalf("expected %d", len(expectedPortList))
	}
	if !reflect.DeepEqual(actualPortList, expectedPortList) {
		t.Fatalf("not equal: \nactual=%v\nexpected=%v", actualPortList, expectedPortList)
	}
}

// Есть два обычных COM-порта.
// Есть один USB COM-порт, VID/PID которого не удовлетворяет условию.
// Есть два USB COM-порта, VID/PID которых удовлетворяют условию.
// Ошибка при открытии порта.
//
// Должен вернуть на фронт ошибку.
func TestUpdatePortList_UsbPortConditionTrue_Open_Fail(t *testing.T) {
	comport := &ComPortMock{}
	app := NewApp(comport)

	comport.EnumerateFunc = func() ([]*enumerator.PortDetails, error) {
		// Все остальные поля - в нули.
		port1 := &enumerator.PortDetails{Name: "COM1"}
		port2 := &enumerator.PortDetails{Name: "COM2"}
		port3 := &enumerator.PortDetails{Name: "COM3", IsUSB: true, VID: "0403", PID: "6001"}
		port4 := &enumerator.PortDetails{Name: "COM4", IsUSB: true, VID: strconv.Itoa(int(0x2e8a)), PID: strconv.Itoa(int(0xf00a))}
		port5 := &enumerator.PortDetails{Name: "COM5", IsUSB: true, VID: strconv.Itoa(int(0x2e8a)), PID: strconv.Itoa(int(0xf00f))}

		var res []*enumerator.PortDetails
		res = append(res, port1, port2, port3, port4, port5)

		return res, nil
	}

	comport.OpenFunc = func(portName string) error {
		return fmt.Errorf("Open failed")
	}

	actualPortList, err := app.UpdatePortList()
	if err == nil {
		t.Fatalf("expected !nil")
	}
	if actualPortList != nil {
		t.Fatalf("expected nil")
	}
}

// Есть два обычных COM-порта.
// Есть один USB COM-порт, VID/PID которого не удовлетворяет условию.
// Есть два USB COM-порта, VID/PID которых удовлетворяют условию.
// Открытие порта.
// Ошибка при записи в порт.
//
// Должен вернуть на фронт ошибку.
func TestUpdatePortList_UsbPortConditionTrue_Write_Fail(t *testing.T) {
	comport := &ComPortMock{}
	app := NewApp(comport)

	comport.EnumerateFunc = func() ([]*enumerator.PortDetails, error) {
		// Все остальные поля - в нули.
		port1 := &enumerator.PortDetails{Name: "COM1"}
		port2 := &enumerator.PortDetails{Name: "COM2"}
		port3 := &enumerator.PortDetails{Name: "COM3", IsUSB: true, VID: "0403", PID: "6001"}
		port4 := &enumerator.PortDetails{Name: "COM4", IsUSB: true, VID: strconv.Itoa(int(0x2e8a)), PID: strconv.Itoa(int(0xf00a))}
		port5 := &enumerator.PortDetails{Name: "COM5", IsUSB: true, VID: strconv.Itoa(int(0x2e8a)), PID: strconv.Itoa(int(0xf00f))}

		var res []*enumerator.PortDetails
		res = append(res, port1, port2, port3, port4, port5)

		return res, nil
	}

	comport.OpenFunc = func(portName string) error { return nil }

	comport.WriteFunc = func(p []byte) (int, error) {
		if !bytes.Equal(p, []byte("AT+VERSION\r\n")) {
			t.Fatalf("expected %s", string([]byte("AT+VERSION\r\n")))
		}
		return 0, fmt.Errorf("Write failed")
	}

	actualPortList, err := app.UpdatePortList()
	if err == nil {
		t.Fatalf("expected !nil")
	}
	if actualPortList != nil {
		t.Fatalf("expected nil")
	}
}

// Есть два обычных COM-порта.
// Есть один USB COM-порт, VID/PID которого не удовлетворяет условию.
// Есть два USB COM-порта, VID/PID которых удовлетворяют условию.
// Открытие порта.
// Запись в порт.
// Ошибка при чтении порта.
//
// Должен вернуть на фронт ошибку.
func TestUpdatePortList_UsbPortConditionTrue_Read_Fail(t *testing.T) {
	comport := &ComPortMock{}
	app := NewApp(comport)

	comport.EnumerateFunc = func() ([]*enumerator.PortDetails, error) {
		// Все остальные поля - в нули.
		port1 := &enumerator.PortDetails{Name: "COM1"}
		port2 := &enumerator.PortDetails{Name: "COM2"}
		port3 := &enumerator.PortDetails{Name: "COM3", IsUSB: true, VID: "0403", PID: "6001"}
		port4 := &enumerator.PortDetails{Name: "COM4", IsUSB: true, VID: strconv.Itoa(int(0x2e8a)), PID: strconv.Itoa(int(0xf00a))}
		port5 := &enumerator.PortDetails{Name: "COM5", IsUSB: true, VID: strconv.Itoa(int(0x2e8a)), PID: strconv.Itoa(int(0xf00f))}

		var res []*enumerator.PortDetails
		res = append(res, port1, port2, port3, port4, port5)

		return res, nil
	}

	comport.OpenFunc = func(portName string) error { return nil }

	comport.WriteFunc = func(p []byte) (int, error) {
		if !bytes.Equal(p, []byte("AT+VERSION\r\n")) {
			t.Fatalf("expected %s", string([]byte("AT+VERSION\r\n")))
		}
		return len([]byte("AT+VERSION\r\n")), nil
	}

	comport.ReadFunc = func(p []byte) (int, error) {
		return 0, fmt.Errorf("Read failed")
	}

	actualPortList, err := app.UpdatePortList()
	if err == nil {
		t.Fatalf("expected !nil")
	}
	if actualPortList != nil {
		t.Fatalf("expected nil")
	}
}

// Есть два обычных COM-порта.
// Есть один USB COM-порт, VID/PID которого не удовлетворяет условию.
// Есть два USB COM-порта, VID/PID которых удовлетворяют условию.
// Открытие порта.
// Запись в порт.
// Чтение порта.
// Ошибка при закрытии порта.
//
// Должен вернуть на фронт ошибку.
func TestUpdatePortList_UsbPortConditionTrue_Close_Fail(t *testing.T) {
	comport := &ComPortMock{}
	app := NewApp(comport)

	comport.EnumerateFunc = func() ([]*enumerator.PortDetails, error) {
		// Все остальные поля - в нули.
		port1 := &enumerator.PortDetails{Name: "COM1"}
		port2 := &enumerator.PortDetails{Name: "COM2"}
		port3 := &enumerator.PortDetails{Name: "COM3", IsUSB: true, VID: "0403", PID: "6001"}
		port4 := &enumerator.PortDetails{Name: "COM4", IsUSB: true, VID: strconv.Itoa(int(0x2e8a)), PID: strconv.Itoa(int(0xf00a))}
		port5 := &enumerator.PortDetails{Name: "COM5", IsUSB: true, VID: strconv.Itoa(int(0x2e8a)), PID: strconv.Itoa(int(0xf00f))}

		var res []*enumerator.PortDetails
		res = append(res, port1, port2, port3, port4, port5)

		return res, nil
	}

	comport.OpenFunc = func(portName string) error { return nil }

	comport.WriteFunc = func(p []byte) (int, error) {
		if !bytes.Equal(p, []byte("AT+VERSION\r\n")) {
			t.Fatalf("expected %s", string([]byte("AT+VERSION\r\n")))
		}
		return len([]byte("AT+VERSION\r\n")), nil
	}

	var deviceNumber = 1
	comport.ReadFunc = func(p []byte) (int, error) {
		str := fmt.Sprintf("v1.2.3 Device%d\r\n", deviceNumber)
		deviceNumber++
		n := copy(p, []byte(str))
		return n, nil
	}

	comport.CloseFunc = func() error {
		return fmt.Errorf("Close failed")
	}

	actualPortList, err := app.UpdatePortList()
	if err == nil {
		t.Fatalf("expected !nil")
	}
	if actualPortList != nil {
		t.Fatalf("expected nil")
	}
}

// Есть два обычных COM-порта.
// Есть один USB COM-порт, VID/PID которого не удовлетворяет условию.
// Есть два USB COM-порта, VID/PID которых удовлетворяют условию.
// Открытие порта.
// Запись в порт.
// Чтение порта.
// Закрытие порта.
//
// Должен вернуть на фронт список всех портов и
// значения Ping/Pong для USB COM-портов, удовлетворяющих условию.
func TestUpdatePortList_UsbPortConditionTrue_Ok(t *testing.T) {
	comport := &ComPortMock{}
	app := NewApp(comport)

	comport.EnumerateFunc = func() ([]*enumerator.PortDetails, error) {
		// Все остальные поля - в нули.
		port1 := &enumerator.PortDetails{Name: "COM1"}
		port2 := &enumerator.PortDetails{Name: "COM2"}
		port3 := &enumerator.PortDetails{Name: "COM3", IsUSB: true, VID: "0403", PID: "6001"}
		port4 := &enumerator.PortDetails{Name: "COM4", IsUSB: true, VID: strconv.Itoa(int(0x2e8a)), PID: strconv.Itoa(int(0xf00a))}
		port5 := &enumerator.PortDetails{Name: "COM5", IsUSB: true, VID: strconv.Itoa(int(0x2e8a)), PID: strconv.Itoa(int(0xf00f))}

		var res []*enumerator.PortDetails
		res = append(res, port1, port2, port3, port4, port5)

		return res, nil
	}

	comport.OpenFunc = func(portName string) error { return nil }

	comport.WriteFunc = func(p []byte) (int, error) {
		if !bytes.Equal(p, []byte("AT+VERSION\r\n")) {
			t.Fatalf("expected %s", string([]byte("AT+VERSION\r\n")))
		}
		return len([]byte("AT+VERSION\r\n")), nil
	}

	var deviceNumber = 1
	comport.ReadFunc = func(p []byte) (int, error) {
		str := fmt.Sprintf("v1.2.3 Device%d\r\n", deviceNumber)
		deviceNumber++
		n := copy(p, []byte(str))
		return n, nil
	}

	comport.CloseFunc = func() error { return nil }

	expectedPortList := []ComPortInfo{
		{Name: "COM1", Usb: "false", Vid: "", Pid: "", SentData: "", ReceivedData: ""},
		{Name: "COM2", Usb: "false", Vid: "", Pid: "", SentData: "", ReceivedData: ""},
		{Name: "COM3", Usb: "true", Vid: "0403", Pid: "6001", SentData: "", ReceivedData: ""},
		{Name: "COM4", Usb: "true", Vid: "11914", Pid: "61450", SentData: "AT+VERSION\r\n", ReceivedData: "v1.2.3 Device1\r\n"},
		{Name: "COM5", Usb: "true", Vid: "11914", Pid: "61455", SentData: "AT+VERSION\r\n", ReceivedData: "v1.2.3 Device2\r\n"},
	}
	actualPortList, err := app.UpdatePortList()
	if err != nil {
		t.Fatalf("expected nil")
	}
	if len(actualPortList) != len(expectedPortList) {
		t.Fatalf("expected %d", len(expectedPortList))
	}
	if !reflect.DeepEqual(actualPortList, expectedPortList) {
		t.Fatalf("not equal: \nactual=%v\nexpected=%v", actualPortList, expectedPortList)
	}
}

// func TestReal(t *testing.T) {
// 	comport := &ComPortReal{}
// 	app := NewApp(comport)

// 	portList, err := app.UpdatePortList()
// 	if err != nil {
// 		t.Fatalf("expected nil")
// 	}
// 	fmt.Println(portList)
// }
