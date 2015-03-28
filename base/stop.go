package main

import (
	"fmt"
	"github.com/ldmberman/GoEV3/Motor"
	"time"
)

func main() {
	fmt.Println("Stop")
	Motor.Stop(Motor.OutPortA)
	Motor.Stop(Motor.OutPortB)
	time.Sleep(time.Second)
}
