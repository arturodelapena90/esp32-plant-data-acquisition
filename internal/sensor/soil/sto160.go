//go:build tinygo

package soil

import (
	"fmt"
	"machine"
	"sync"
)

var initADCOnce sync.Once

func initSoilADC(pin machine.Pin) (machine.ADC, error) {
	initADCOnce.Do(machine.InitADC)

	adc := machine.ADC{Pin: pin}
	adc.Configure(machine.ADCConfig{})

	return adc, nil
}

func readSoilADC(adc machine.ADC) (*float32, error) {
	raw := uint32(adc.Get())

	// Get() returns a 16-bit-scaled value (0..65520, per machine_esp32s3_adc.go),
	percentage := float32(raw) / 65520 * 100

	fmt.Printf(
		"soil ADC reading: raw=%d moisture=%.2f%%\n",
		raw,
		percentage,
	)

	return &percentage, nil
}
