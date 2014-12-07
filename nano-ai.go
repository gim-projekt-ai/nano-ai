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

type SynonymePair struct {
	base string
	synonyme string
}

func main() {
	fmt.Println("nano-ai 0.0.1")
	//zmienne
	//arr := []string{"g","k"}
	//AddNotSynonyme(arr)
	//fmt.Printf("%v", Scannsyn())
	
	var unprocQuery string
	var purpose int8
	//główna pętla
	for {
		unprocQuery = GetQuery()
		fmt.Println("------------------Informacje--------------------")
		if unprocQuery == "*quit" {
			os.Exit(1)
		}
		purpose = Querypurpose(unprocQuery)
		//typ qypowiedzi: 1 - inform.; 2- pyt. 3 - rozkaz
		println(purpose, unprocQuery)
		if purpose == 1 {
			addtodb(unprocQuery)
		}
		if purpose == 2 {
			dbcontents := Scandb()
			//fmt.Printf("%d\n", dbcontents)
			placeofq := strings.Index(unprocQuery, "*q")
			qprefix := unprocQuery[:placeofq]
			qsuffix := unprocQuery[(placeofq + 2):]
			//fmt.Print(string(placeofq), qprefix,"  ,  ", qsuffix)
			response := GrepIn(dbcontents, qprefix, qsuffix)
			fmt.Println("-----------------Info-koniec--------------------")
			fmt.Printf("%s\n", response[0])

		}
		//synonimy

	}
}

func GrepIn(contents []string, qprefix, qsuffix string) []string {
	itHasPrefix := make([]string, 160)
	answers := make([]string, 320)
	var pcount int8 = 0
	var acount int8 = 0

	for _, v := range contents {
		if strings.HasPrefix(v, qprefix) {
			//naprawienie błędu z pustymi możliwościami
			if strings.Trim(v, " ") != "" {
				fmt.Printf("GrepIn Prefix: %s\n", v)
				itHasPrefix[pcount] = v
				pcount = pcount + 1
			}
		}
	}
	for _, v := range itHasPrefix {
		if strings.HasSuffix(v, qsuffix) {
			//naprawienie błędu
			if strings.Trim(v, " ") != "" {
				fmt.Printf("GrepIn Suffix: %s\n", v)
				answers[acount] = v
				acount = acount + 1
			}
		}
	}
	//fmt.Printf("%d\n", answers)
	return answers
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
func Scandb() []string {
	dat, err := ioutil.ReadFile("db1.txt")
	errorcheck(err)
	data := string(dat)
	//informatycznie odpowiednia długość wycinka
	lines := make([]string, 16383)
	lines = strings.Split(data, "\n")
	return lines
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
	// jako odczyt/zapis, dopisywanie, tworzy nowy plik jeżeli nie ma.
	f, err := os.OpenFile("db1.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	errorcheck(err)
	//zakodowanie linii
	d2 := []byte(query + "\n")
	//piszemy
	n2, err := f.Write(d2)
	errorcheck(err)
	defer fmt.Printf("wrote %d bytes\n", n2)
}

func Scansyn() []string {
	dat, err := ioutil.ReadFile("db1.txt")
	errorcheck(err)
	data := string(dat)
	//informatycznie odpowiednia długość wycinka
	lines := make([]string, 16383)
	lines = strings.Split(data, "\n")
	return lines
}
func Scannsyn() []string {
	dat, err := ioutil.ReadFile("nsyn.txt")
	errorcheck(err)
	data := string(dat)
	//informatycznie odpowiednia długość wycinka
	lines := make([]string, 16383)
	lines = strings.Split(data, "\n")
	return lines
}
func AddNotSynonyme(pair []string) {
	f, err := os.OpenFile("nsyn.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	errorcheck(err)
	//zakodowanie linii
	d2 := []byte(pair[0] + " " + pair[1] + "\n")
	//piszemy
	n2, err := f.Write(d2)
	errorcheck(err)
	defer fmt.Printf("wrote %d bytes\n", n2)
}
func AddSynonyme(pair []string) {
	f, err := os.OpenFile("syn.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	errorcheck(err)
	//zakodowanie linii
	d2 := []byte(pair[0] + " " + pair[1] + "\n")
	//piszemy
	n2, err := f.Write(d2)
	errorcheck(err)
	defer fmt.Printf("wrote %d bytes\n", n2)
}
func BaseWordOf(word string) string {
	synonymes := Scansyn()
	
}

/*
func FindAnalogical(query SlicedQuery) []string {
	abvc
}

func RemoveSynonymes(query SlicedQuery) SlicedQuery {
	abc
}

func GrepForSynonymes {
	abc
}
type SlicedQuery struct {
	type int8
	obj, verb, subject string
	logic string
	obj2, verb2, subject2 string
}


*/

func errorcheck(e error) {
	if e != nil {
		panic(e)
	}
}
