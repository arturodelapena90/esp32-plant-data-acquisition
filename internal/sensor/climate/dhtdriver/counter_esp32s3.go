//go:build tinygo && esp32s3

package dhtdriver

// esp32s3 runs at up to 240MHz, well above the 2^16 ticks-per-millisecond
// range a uint16 counter can hold (upstream's highfreq.go doesn't list
// esp32s3, so without this it would silently wrap on the lowfreq uint16
// path). See doc.go for why this package is vendored.
type counter uint32
