//go:build tinygo

package soil

import (
	"machine"
	"sync"
)

var initADCOnce sync.Once

// calibration
const (
	dryRawPercent = 60.7
	wetRawPercent = 23.8
)

func initSoilADC(pin machine.Pin) (machine.ADC, error) {
	initADCOnce.Do(machine.InitADC)

	adc := machine.ADC{Pin: pin}
	adc.Configure(machine.ADCConfig{})

	return adc, nil
}

func readSoilADC(adc machine.ADC) (*float32, error) {
	raw := uint32(adc.Get())

	// Get() returns a 16-bit-scaled value (0..65520, per machine_esp32s3_adc.go).
	rawPercent := float32(raw) / 65520 * 100

	moisture := (dryRawPercent - rawPercent) / (dryRawPercent - wetRawPercent) * 100
	if moisture < 0 {
		moisture = 0
	} else if moisture > 100 {
		moisture = 100
	}

	return &moisture, nil
}
