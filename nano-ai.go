/*MIT licence
(c) Jan Piskorski, Bartosz Deptuła, Mateusz Boruń
Wiktor Kacalski & Mikołaj Kordowski
*/

package main

import (
	"fmt"
//	"strings"
	"bufio"
	"os"
)

func main() {
	fmt.Println("nano-ai 0.0.1")
	var unprocQuery string
	var purpose int8
	for {
		unprocQuery = GetQuery()
		purpose = Querypurpose(unprocQuery)
		//typ qypowiedzi: 1 - inform.; 2- pyt. 3 - rozkaz
		println(purpose, unprocQuery)
		os.Exit(0)
	}

}
func GetQuery() string {
	var inp string
	scnr := bufio.NewScanner(os.Stdin)
	scnr.Scan()
	inp = scnr.Text()
	//fmt.Printf("%s\n", scnr.Text())
	return inp
}
func Querypurpose(query string) int8 {
	return 1
}
