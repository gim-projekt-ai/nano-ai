/*MIT licence
(c) Jan Piskorski, Bartosz Deptuła, Mateusz Boruń
Wiktor Kacalski & Mikołaj Kordowski
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("nano-ai 0.0.1")
	//zmienne
	var unprocQuery string
	var purpose int8
	//główna pętla
	for {
		unprocQuery = GetQuery()
		purpose = Querypurpose(unprocQuery)
		//typ qypowiedzi: 1 - inform.; 2- pyt. 3 - rozkaz
		println(purpose, unprocQuery)
		//na razie wychodzi
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
	var querytype int8 = 1
	//jak zawiera *q to jest pytaniem
	if strings.Contains(query, "*q") {
		querytype = 2
	}
	//jak zawiera request to jest rozkazem
	if strings.Contains(query, "*request") {
		querytype = 3
	}
	//jeśliu nie zawiera żadnego z pow. to nadal =1
	return querytype

}
