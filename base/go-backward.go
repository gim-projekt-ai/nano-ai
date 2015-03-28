package main

import (
	"fmt"
	"github.com/ldmberman/GoEV3/Motor"
	"time"
)

func main() {
	fmt.Println("Backward")
	Motor.Run(Motor.OutPortA, -50)
	Motor.Run(Motor.OutPortB, -50)
	time.Sleep(time.Second*3)
	/*
	for {
		if Motor.CurrentSpeed(Motor.OutPortA) > -10 {
			Motor.Stop(Motor.OutPortA)
			Motor.Stop(Motor.OutPortB)
			break
		}
		if Motor.CurrentSpeed(Motor.OutPortB) > -10 {
			Motor.Stop(Motor.OutPortA)
			Motor.Stop(Motor.OutPortB)
			break
		}
		time.Sleep(time.Second)
	}*/
	Motor.Stop(Motor.OutPortA)
	Motor.Stop(Motor.OutPortB)
}
