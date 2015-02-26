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
var verbose bool = false

func main() {
	fmt.Println("nano-ai 0.2.1")
	t := time.Now
	log("Nano-AI log: ", t.Local())
	//troche przygotowan
	fmt.Print("Scanning for unprocessed synonymes...")
	SynonymeCheck()
	fmt.Println("done")
	fmt.Print("Removing doubled information...")
	DoubledInfoCheck()
	fmt.Println("done")

	//Procedury testowe
	//brak
	
	verbose = YesNoQuestion("Do you want verbose information? ")
	//zmienne
	var unprocQuery string
	var purpose int8
	//główna pętla
	for {
		unprocQuery = GetQuery()

		vout("------------------Informacje--------------------")
		if unprocQuery == "*quit" {
			os.Exit(1)
		}
		purpose = Querypurpose(unprocQuery)
		//typ qypowiedzi: 1 - inform.; 2- pyt. 3 - rozkaz
		vout("purpose:  ",purpose, unprocQuery)
		dbcontents := Scandb() /*skanowanie db1 */

		if purpose == 1 {
			unprocQuery = RemoveSynonymes(unprocQuery)
			//Trafne spostrzezenie dot. informacji
			fmt.Println(Type1Response(dbcontents, unprocQuery))
			//synonimy i zapamietanie
			SynonymeManagement(unprocQuery)
			addtodb(unprocQuery)
		}
		if purpose == 2 {
			vout(dbcontents,"\n")
			placeofq := strings.Index(unprocQuery, "*q")
			qprefix := unprocQuery[:placeofq]
			qsuffix := unprocQuery[(placeofq + 2):]

			qprefix = strings.Trim(qprefix, " _\t\n!.")
			qsuffix = strings.Trim(qsuffix, " _\t\n!.")
			vout(string(placeofq), qprefix,"  ,  ", qsuffix)
			response := GrepIn(dbcontents, qprefix, qsuffix)
			vout("-----------------Info-koniec-----------------")
			fmt.Println(Type2Response(qprefix, qsuffix, response))

		}
		//synonimy

	}
}

/* Wyszukiwanie odpowiedzi z zadanym początkiem i końcem. Sprawia problemy, ale tylko czasem.
 * Daje się okiełznać.
 */
func GrepIn(contents []string, qprefix, qsuffix string) []string {
	itHasPrefix := make([]string, 160)
	answers := make([]string, 160)
	var pcount int8 = 0
	var acount int8 = 0
	vout("qprefix:", qprefix ,"\nqsuffix:", qsuffix ,"\n")

	for _, v := range contents {
		if strings.HasPrefix(v, qprefix) {
			//naprawienie błędu z pustymi możliwościami
			if strings.Trim(v, " ") != "" {
				//fmt.Printf("GrepIn Prefix: %s\n", v)
				itHasPrefix[pcount] = strings.Trim(v, " \t\n")
				pcount = pcount + 1
			}
		}
	}
	//fmt.Printf("GrepInPrefixMatch: %d\n", itHasPrefix)
	for _, v := range itHasPrefix {
		if strings.HasSuffix(v, qsuffix) {
			//naprawienie błędu
			if strings.Trim(v, " ") != "" {
				//fmt.Printf("GrepIn Suffix: %s\n", v)
				answers[acount] = v
				acount = acount + 1
			}
		}
	}
	//Nadal puste?
	/*
		answ1 := make([]string, 160)
		var ecount int8 = 0
		for _, v:=range answers {
			if (strings.Trim(v," \t\n_")) != "" {
				answ1[ecount] = v
				ecount++
			}
		}
	*/
	vout("GrepInMatch:",answers, "\n")
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
func YesNoQuestion(q string) bool {
	var inp string
	fmt.Printf("%v ", q)
	scnr := bufio.NewScanner(os.Stdin)
	scnr.Scan()
	inp = scnr.Text()
	var o bool
	if ((inp[:1] == "t") || (inp[:1] == "y")) || ((inp[:1] == "Y") || (inp[:1] == "T")) {
		o = true
	} else {
		o = false
	}
	return o
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
func SliceAndTrim(query string) []string {
	sliced := strings.Split(query, " ")
	trimmed := make([]string, 7)
	if query != "" {
		var v string
		for i := 0; i < 3; i++ {
			v = sliced[i]
			trimmed[i] = v[:strings.Index(v, "(")]
		}
	}
	return trimmed
}
func Type1Response(db []string, query string) string {
	response := "Ok, didn't know that."
	sQuery := SliceAndTrim(RemoveSynonymes(query))
	for _, v := range db {
		v1 := SliceAndTrim(v)
		if v1[0] == sQuery[0] {
			response = "And I suppose it " + v1[1] + " " + v1[2] + " also..."
		}
		if (v1[1] == sQuery[1]) && (v1[2] == sQuery[2]) {
			response = "It's true for " + v1[0] + " too!"
			if v1[0] == sQuery[0] {
				response = response + "\nDon't think I'm stupid. I already know that."
			}
		}
	}
	return response
}
func Type2Response(qp, qs string, matching []string) string {
	response := ""
	//qp, qs = strings.Trim(qp, " "), strings.Trim(qs, " ")
	if (qp == "") && (qs == "") {
		if YesNoQuestion("Really print all I know? ") {
			response = strings.Join(matching, ", \n")
		} else {
			response = "Ok."
		}
	} else if qs == "" {
		response = qp
		for _, v := range matching {
			if strings.Trim(v, " ") != "" {
				v1 := SliceAndTrim(v)
				response = response + " " + v1[1] + " " + v1[2] + ",\n"
			}
		}
	} else if qp == "" {
		for _, v := range matching {
			if strings.Trim(v, " ") != "" {
				v1 := SliceAndTrim(v)
				response = response + " " + v1[0] + " " + v1[1] + ",\n"
			}
		}
		response = response + qs
		//response= strings.Join(matching, ", ")
	} else {
		response = qp
		for _, v := range matching {
			if strings.Trim(v, " ") != "" {
				v1 := SliceAndTrim(v)
				response = response + " " + v1[1] + ", "
			}
		}
		response = response + qs
		//response = matching[0]
	}
	return response
}
func DoubledInfoCheck() {
	db :=Scandb()
	undoubleddb := make([]string, len(db))
	count := 0
	dblcnt := 0
	for _, v:= range db {
		if !(in(undoubleddb, strings.Trim(v, " \n\t"))) {
			undoubleddb[count] = strings.Trim(v, " \t\n")
			count += 1
		} else {
			dblcnt += 1
		}
	}
	f, err := os.Create("db1.txt")
	defer f.Close()
	errorcheck(err)
	for _, v := range undoubleddb {
		if v != "" {
			d2 := []byte(RemoveSynonymes(v) + "\n")
			_, err := f.Write(d2)
			errorcheck(err)
			//fmt.Printf("wrote %d bytes\n", n2)
		}
	}
	fmt.Printf("%v found...", (dblcnt-1))
}

func Scansyn() []string {
	dat, err := ioutil.ReadFile("syn.txt")
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
	line := make([]string, 2)
	for _, v := range synonymes {
		line = strings.Split(strings.Trim(v, " "), " ")
		if line[0] != "" {
			if line[1] == word {
				return line[0]
			}
			if line[0] == word {
				return line[0]
			}
		}
		line = make([]string, 2)
	}
	return word
}
func RemoveSynonymes(query string) string {
	sliced := strings.Split(query, " ")
	removed := make([]string, 7)
	var v string
	for i := 0; i < 3; i++ {
		v = sliced[i]
		//fmt.Println(i)
		removed[i] = BaseWordOf(v[:strings.Index(v, "(")]) + v[strings.Index(v, "("):]
		//fmt.Println(removed)
	}
	//fmt.Println(sliced)
	return strings.Join(removed, " ")
}
func PotentialSynonyme(query1, query2 []string) (bool, string, string) {
	var s1, s2 string
	var rly bool = false
	s1 = ""
	s2 = ""
	if query1[0] == query2[0] {
		if query1[1] == query2[1] {
			s1 = query1[2]
			s2 = query2[2]
			rly = true
		} else if query1[2] == query2[2] {
			s1 = query1[1]
			s2 = query2[1]
			rly = true
		}
	}
	if (query1[1] == query2[1]) && (query1[2] == query2[2]) {
		s1 = query1[0]
		s2 = query2[0]
		rly = true
	}
	for _, v := range query1 {
		for _, v1 := range query2 {
			if !((v == "") || (v1 == "")) {
				if (len(v) >= 2) && (len(v1) >= 2) {
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
			}
		}
	}
	return rly, s1, s2
}
func SynonymeManagement(query string) {
	allnsyn := Scannsyn()
	allsyn := Scansyn()
	dbcontents := Scandb()
	for _, v := range dbcontents {
		//fmt.Printf("%v; %v\n", v, query)
		issyn, base, syn := PotentialSynonyme(SliceAndTrim(v), SliceAndTrim(query))
		if issyn {
			if !(strings.Contains(strings.Join(allnsyn, ","), syn)) {
				if !(strings.Contains(strings.Join(allnsyn, ","), base)) {
					if !(strings.Contains(strings.Join(allsyn, ","), syn)) {
						if syn != base {
							if YesNoQuestion("Is " + base + " same as " + syn + "?") {
								AddSynonyme([]string{base, syn})
							} else {
								AddNotSynonyme([]string{base, syn})
							}
						}
					}
				}
			}
		}
	}
}
func SynonymeCheck() {
	db := Scandb()
	for _, v := range db {
		SynonymeManagement(v)
	}
	f, err := os.Create("db1.txt")
	defer f.Close()
	errorcheck(err)
	for _, v := range db {
		if v != "" {
			d2 := []byte(RemoveSynonymes(v) + "\n")
			_, err := f.Write(d2)
			errorcheck(err)
			//fmt.Printf("wrote %d bytes\n", n2)
		}
	}
}



func errorcheck(e error) {
	if e != nil {
		panic(e)
	}
}
func in(sl []string, s string) bool {
	var rlyin bool = false
	for _, v := range sl {
		if v == s {
			rlyin = true
		}
	}
	return rlyin
}
func vout(a ...interface{}) {
	if verbose {
		fmt.Println(a...)
	}
}
func log(a ...interface{}) {
	s = fmt.Sprintln(a...)
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	errorcheck(err)
	d2 := []byte(s + "\n")
	f.Write(d2)
	errorcheck(err)
}
