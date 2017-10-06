package main

import "fmt"
import "github.com/mdlayher/lmsensors"

func main(){
	s:=lmsensors.New()
	ds,err:=s.Scan()
	if err!=nil{
		panic(err)
	}
	for _,d := range ds{
		for _,s := range d.Sensors{
			switch ts:=s.(type){
			case *lmsensors.FanSensor: 
				print(d.Name,*ts)
			case *lmsensors.CurrentSensor: 
				print(d.Name,*ts)
			case *lmsensors.TemperatureSensor: 
				print(d.Name,*ts)
			case *lmsensors.PowerSensor: 
				print(d.Name,*ts)
			case *lmsensors.VoltageSensor: 
				print(d.Name,*ts)
			case *lmsensors.IntrusionSensor: 
				print(d.Name,*ts)

			}		
		}
	}
}

func print(s string, i interface{}){
	fmt.Printf("Device Name:'%s'  Sensor:%[2]T%+[2]v\n",s, i)	
}

