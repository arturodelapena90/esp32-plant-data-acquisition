//go:build tinygo

package dhtdriver

import (
	"time"
)

// Celsius and Fahrenheit temperature scales
type TemperatureScale uint8

func (t TemperatureScale) convertToFloat(temp int16) float32 {
	if t == C {
		return float32(temp) / 10
	} else {
		// Fahrenheit
		return float32(temp)*(9.0/50.) + 32.
	}
}

// All functions return ErrorCode instance as error. This class can be used for more efficient error processing
type ErrorCode uint8

const (
	startTimeout = time.Millisecond * 200
	startingLow  = time.Millisecond * 20

	C TemperatureScale = iota
	F

	ChecksumError ErrorCode = iota
	NoSignalError
	NoDataError
	UpdateError
	UninitializedDataError
)

// error interface implementation for ErrorCode
func (e ErrorCode) Error() string {
	switch e {
	case ChecksumError:
		// DHT returns ChecksumError if all the data from the sensor was received, but the checksum does not match.
		return "checksum mismatch"
	case NoSignalError:
		// DHT returns NoSignalError if there was no reply from the sensor. Check sensor connection or the correct pin
		// sis chosen,
		return "no signal"
	case NoDataError:
		// DHT returns NoDataError if the connection was successfully initialized, but not all 40 bits from
		// the sensor is received
		return "no data"
	case UpdateError:
		// DHT returns UpdateError if ReadMeasurements function is called before time specified in UpdatePolicy or
		// less than 2 seconds after past measurement
		return "cannot update now"
	case UninitializedDataError:
		// DHT returns UninitializedDataError if user attempts to access data before first measurement
		return "no measurements done"
	}
	// should never be reached
	return "unknown error"
}

// Update policy of the DHT device. UpdateTime cannot be shorter than 2 seconds. According to dht specification sensor
// will return undefined data if update requested less than 2 seconds before last usage
type UpdatePolicy struct {
	UpdateTime          time.Duration
	UpdateAutomatically bool
}

var (
	// timeout counter equal to number of ticks per 1 millisecond
	timeout counter
)

func init() {
	timeout = cyclesPerMillisecond()
}

// cyclesPerMillisecond and the "counter" type it returns are defined per
// chip family in freq_esp32s3.go / freq_default.go and
// counter_esp32s3.go / highfreq.go / lowfreq.go.
