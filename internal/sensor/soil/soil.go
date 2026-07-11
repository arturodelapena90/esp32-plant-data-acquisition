package soil

import (
	"fmt"
	"machine"
	"time"
)

type Sensor struct {
	adc machine.ADC
}

type Reading struct {
	Moisture *float32
}

func New(pin machine.Pin) (*Sensor, error) {
	adc, err := initSoilADC(pin)
	if err != nil {
		return nil, err
	}

	return &Sensor{
		adc: adc,
	}, nil
}

func (s *Sensor) Start(interval time.Duration, readingChan chan<- Reading) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		reading, err := s.Read()
		if err != nil {
			fmt.Printf("soil sensor error: %v\n", err)
		}

		readingChan <- reading
	}
}

func (s *Sensor) Read() (Reading, error) {
	moisture, err := readSoilADC(s.adc)

	return Reading{
		Moisture: moisture,
	}, err
}
