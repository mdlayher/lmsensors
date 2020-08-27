package main

import (
	"log"

	"github.com/mdlayher/lmsensors"
)

func main() {
	scanner := lmsensors.New()

	devices, err := scanner.Scan()

	if err != nil {
		log.Fatal(err)
	}

	for _, d := range devices {
		for _, s := range d.Sensors {
			handleSensors(s)
		}
	}
}

func handleSensors(i interface{}) {
	switch v := i.(type) {
	case *lmsensors.CurrentSensor:
		log.Println(v)
	case *lmsensors.VoltageSensor:
		log.Println(v)
	case *lmsensors.FanSensor:
		log.Println(v)
	case *lmsensors.IntrusionSensor:
		log.Println(v)
	case *lmsensors.PowerSensor:
		log.Println(v)
	case *lmsensors.TemperatureSensor:
		log.Println(v)
	default:
		log.Println("Undetected Sensor Type")
	}
}
