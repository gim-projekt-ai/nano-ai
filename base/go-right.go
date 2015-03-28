package main

import (
	"fmt"
	"github.com/ldmberman/GoEV3/Motor"
	"time"
)

func main() {
	fmt.Println("Right")
	go Motor.Run(Motor.OutPortA, -50)
	go Motor.Run(Motor.OutPortB, 50)
	time.Sleep(time.Millisecond*1300)
	/*
	for {
		if Motor.CurrentSpeed(Motor.OutPortA) > -10 {
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
	}
	*/
	Motor.Stop(Motor.OutPortA)
	Motor.Stop(Motor.OutPortB)
}
