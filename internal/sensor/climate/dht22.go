//go:build tinygo

package climate

import (
	"fmt"
	"machine"

	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/climate/dhtdriver"
)

func initDHT22(pin machine.Pin) (dhtdriver.Device, error) {
	device := dhtdriver.New(pin, dhtdriver.DHT22)
	return device, nil
}

func readDHT22(device dhtdriver.Device) (*float32, *float32, error) {
	if err := device.ReadMeasurements(); err != nil {
		fmt.Printf("DHT22 read error: %v\n", err)
		return nil, nil, err
	}

	temp, humidity, err := device.Measurements()
	if err != nil {
		fmt.Printf("DHT22 measurements error: %v\n", err)
		return nil, nil, err
	}

	tempFloat := float32(temp) / 10.0
	humiFloat := float32(humidity) / 10.0

	fmt.Printf(
		"climate reading: %.1f°C, %.1f%% humidity\n",
		tempFloat,
		humiFloat,
	)

	return &tempFloat, &humiFloat, nil
}
