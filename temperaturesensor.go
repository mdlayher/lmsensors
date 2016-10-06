package lmsensors

import (
	"strconv"
)

var _ Sensor = &TemperatureSensor{}

// A TemperatureSensor is a Sensor that detects temperatures in degrees
// Celsius.
type TemperatureSensor struct {
	// The name of the sensor.
	Name string

	// A label that describes what the sensor is monitoring.  Label may be
	// empty.
	Label string

	// The current temperature, in degrees Celsius, indicated by the sensor.
	Current float64

	// A high threshold temperature, in degrees Celsius, indicated by the
	// sensor.
	High float64

	// A critical threshold temperature, in degrees Celsius, indicated by the
	// sensor.
	Critical float64

	// Whether or not the temperature is past the critical threshold.
	CriticalAlarm bool
}

func (s *TemperatureSensor) name() string        { return s.Name }
func (s *TemperatureSensor) setName(name string) { s.Name = name }

func (s *TemperatureSensor) parse(raw map[string]string) error {
	for k, v := range raw {
		switch k {
		case "input", "crit", "max":
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}

			// Raw temperature values are scaled by 1000
			f /= 1000

			switch k {
			case "input":
				s.Current = f
			case "crit":
				s.Critical = f
			case "max":
				s.High = f
			}
		case "crit_alarm":
			s.CriticalAlarm = v != "0"
		case "label":
			s.Label = v
		}
	}

	return nil
}
