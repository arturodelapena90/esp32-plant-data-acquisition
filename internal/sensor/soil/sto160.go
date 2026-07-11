//go:build tinygo

package soil

import (
	"fmt"
	"machine"
)

func initSoilADC(pin machine.Pin) (machine.ADC, error) {
	adc := machine.ADC{Pin: pin}
	adc.Configure(machine.ADCConfig{})

	return adc, nil
}

func readSoilADC(adc machine.ADC) (*float32, error) {
	raw := uint32(adc.Get())

	percentage := float32(raw) / 4095 * 100

	fmt.Printf(
		"soil ADC reading: raw=%d moisture=%.2f%%\n",
		raw,
		percentage,
	)

	return &percentage, nil
}
