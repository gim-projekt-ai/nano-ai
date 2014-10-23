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
	//fmt.Println("nano 0.0.1")
	fmt.Println("Napisaleś: "+GetQuery())

}
func GetQuery() string {
	var inp string
	scnr := bufio.NewScanner(os.Stdin)
	scnr.Scan()
	inp = scnr.Text()
	//fmt.Printf("%s\n", scnr.Text())
	return inp
}
