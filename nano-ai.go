/*MIT licence
(c) Jan Piskorski, Bartosz Deptuła, Mateusz Boruń
Wiktor Kacalski & Mikołaj Kordowski
*/

package main

import (
	//print
	"fmt"
	//pobieranie danych
	"bufio"
	"os"
	//operacje na słowach
	"strings"
	//pliki - na razie wystarczy os
	"io/ioutil"
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
		if purpose == 1 {
			addtodb(unprocQuery)
		}
		if purpose == 2 {
			fmt.Print(Scandb(unprocQuery) + "\n")
		}
		//if purpose == 2 {
		//na razie wychodzi
		//os.Exit(0)
	}

}
func GetQuery() string {
	var inp string
	fmt.Printf(";> ")
	//źródło to konsola
	scnr := bufio.NewScanner(os.Stdin)
	//skanujemy i wynik do zmiennej
	scnr.Scan()
	inp = scnr.Text()
	//fmt.Printf("%s\n", scnr.Text())
	return inp
}
func Scandb(content string) string {
	dat, err := ioutil.ReadFile("db1.txt")
	errorcheck(err)
	data := string(dat)
	return data
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
	//jeśli nie zawiera żadnego z pow. to nadal =1
	return querytype

}
func addtodb(query string) {
	/*f, err := os.Create("db1.txt")
	  errorcheck(err)
	  d2 := []byte(query+"\n")
	  n2, err := f.Write(d2)
	  errorcheck(err)
	  fmt.Printf("wrote %d bytes\n", n2)
	  f.Close()*/
	// jako odczyt/zapis, dopisywanie, tworzy nowy jeżeli nie ma.
	f, err := os.OpenFile("db1.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	errorcheck(err)
	//zakodowanie linii
	d2 := []byte(query + "\n")
	//piszemy
	n2, err := f.Write(d2)
	errorcheck(err)
	defer fmt.Printf("wrote %d bytes\n", n2)
	f.Close()
}

func errorcheck(e error) {
	if e != nil {
		panic(e)
	}
}
