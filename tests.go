package main

import (
	"fmt"
)

func main() {
	var s1, s2 string
	var rly bool = false
	s1 = ""
	s2 = ""
	query1 := []string{"ala","ma","kota"}
	query2 := []string{"ala", "masz", "kota"}
	for _, v:= range query1 {
		for _, v1:= range query2 {
			
			if !((v == "") || (v1 == "")) {
				if (v == v1[:len(v1)-1]) || (v == v1[:len(v1)-2]) {
					s1 = v
					s2 = v1
					rly = true
				}
				if (v1 == v[:len(v)-1]) || (v1 == v[:len(v)-2]) {
					s1 = v1
					s2 = v
					rly = true
				}
			}
			fmt.Printf("%v   ",v1)
		}
		fmt.Printf("\n+%v\n", v)
	}
	fmt.Printf("%v;%v;%v; \n", rly, s1, s2)
}
