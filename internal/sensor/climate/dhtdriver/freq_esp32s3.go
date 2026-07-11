//go:build tinygo && esp32s3

package dhtdriver

import "machine"

// esp32s3 has no machine.CPUFrequency(); it only exposes
// GetCPUFrequency() (uint32, error). See doc.go for why this package is
// vendored.
func cyclesPerMillisecond() counter {
	freq, err := machine.GetCPUFrequency()
	if err != nil {
		// Default esp32s3 CPU clock per the datasheet; used only if the
		// runtime query fails.
		freq = 240_000_000
	}
	freq /= 1000
	return counter(freq)
}
