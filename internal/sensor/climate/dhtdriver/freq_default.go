//go:build tinygo && !esp32s3

package dhtdriver

import "machine"

func cyclesPerMillisecond() counter {
	freq := machine.CPUFrequency()
	freq /= 1000
	return counter(freq)
}
