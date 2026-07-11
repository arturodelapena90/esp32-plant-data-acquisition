// Package dhtdriver is a vendored copy of tinygo.org/x/drivers/dht@v0.35.0,
// under the BSD-3-Clause license in LICENSE (Copyright The TinyGo Authors).
//
// It is vendored rather than imported because that release (and the "dev"
// branch as of 2026-07) fails to build for esp32s3: it calls
// machine.CPUFrequency(), which esp32s3 does not implement (only
// machine.GetCPUFrequency() (uint32, error) exists there), and even if it
// did, esp32s3's ~240MHz clock would overflow the uint16 "counter" type
// picked for chips assumed to run below 64MHz. See freq_esp32s3.go and
// counter_esp32s3.go for the chip-specific fix; freq_default.go and the
// upstream highfreq.go/lowfreq.go are unchanged for all other chips.
package dhtdriver // import "github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/climate/dhtdriver"
