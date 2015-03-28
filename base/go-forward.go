package main

import (
	"fmt"
	"github.com/ldmberman/GoEV3/Motor"
	"time"
	//"tinyev3lib"
)

func main() {
	fmt.Println("Forward")
	//var tim string
	//_, _ = fmt.Scanf("%s", &tim)
	go Motor.Run(Motor.OutPortA, 75)
	go Motor.Run(Motor.OutPortB, 75)
	//tim2, _ := time.ParseDuration(tim+"s")
	time.Sleep(time.Second*2)
	/*
	for {
		if Motor.CurrentSpeed(Motor.OutPortA) < 10 {
			Motor.Stop(Motor.OutPortA)
			Motor.Stop(Motor.OutPortB)
			break
		}
		if Motor.CurrentSpeed(Motor.OutPortB) < 10 {
			Motor.Stop(Motor.OutPortA)
			Motor.Stop(Motor.OutPortB)
			break
		}
		time.Sleep(time.Second)
	}*/
	Motor.Stop(Motor.OutPortA)
	Motor.Stop(Motor.OutPortB)
}
 

 
