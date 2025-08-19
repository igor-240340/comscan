package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	comport := &ComPortReal{}

	// NOTE: Для простоты передаем моку и тестируем фронт вручную.
	// comport := &ComPortMock{}
	// BEGIN: Успешный сценарий.
	/*
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
	*/
	// END: Успешный сценарий.

	// BEGIN: Ошибка при вызове Enumerate().
	/*
		comport.EnumerateFunc = func() ([]*enumerator.PortDetails, error) {
			return nil, fmt.Errorf("Enumerate failed")
		}
	*/
	// END: Ошибка при вызове Enumerate().

	app := NewApp(comport)

	err := wails.Run(&options.App{
		Title:         "comscan",
		Width:         800,
		Height:        600,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
