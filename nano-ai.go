/*MIT licence
(c) Jan Piskorski
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
	"time"
	//czytanie classif
	"path/filepath"
	//NLP - niestandardowa
	"github.com/gim-projekt-ai/nanoai-libs/convert_new"
	
	//Zadania - niestandardowa
	"github.com/gim-projekt-ai/nanoai-libs/aiRequests"
	//losowanie
	"math/rand"
)

var verbose bool = false
var textformatting bool = false
var onRobot bool = false

func main() {
	fmt.Println("nano-ai 0.3.2")
	log("Witaj w logu Nano-AI! 0.3.2")
	//troche przygotowan
	fmt.Print("Scanning for unprocessed synonymes...")
	SynonymeCheck()
	fmt.Println("done")
	fmt.Print("Removing doubled information...")
	DoubledInfoCheck()
	fmt.Println("done")
	fmt.Print("Removing doubled classification...")
	DoubledClassificationCheck()
	fmt.Println("found...done")
	log("Zakończono przygotowania.")

	//Procedury testowe
	fmt.Println(ListClassification())

	verbose = YesNoQuestion("Do you want verbose information? ")
	textformatting = YesNoQuestion("Do you want the text to be NLP-ed?")
	onRobot = YesNoQuestion("Is the program running on a robot?")
	//zmienne
	var unprocQuery string
	var purpose int8
	//główna pętla
	for {
		unprocQuery = GetQuery()

		vout("------------------Informacje--------------------")
		if unprocQuery == "*quit" {
			log("Do widzenia!")
			os.Exit(0)
		}
		if unprocQuery == "*empty" {
			log("Puste zdanie!")
			fmt.Println("Please coop.")
			continue
		}
		purpose = Querypurpose(unprocQuery)
		//typ qypowiedzi: 1 - inform.; 2- pyt. 3 - rozkaz
		vout("purpose:  ", purpose, unprocQuery)
		dbcontents := Scandb() /*skanowanie db1 */
		
		if purpose == 1 {
			log("Pobrano ", unprocQuery, ", typu informacja.")
			unprocQuery = RemoveSynonymes(unprocQuery)
			//Trafne spostrzezenie dot. informacji
			fmt.Println(Type1Response(dbcontents, unprocQuery))
			//synonimy i zapamietanie
			SynonymeManagement(unprocQuery)
			addtodb(unprocQuery)
		}
		if purpose == 101 {
			welcomes := []string{"Hi, nice to meet you!", "Hello!",
				"Welcome!", "Hi!", "Guten Tag!", "Good Morning!"}
			fmt.Println(welcomes[rand.Intn(len(welcomes))])
		}
		if purpose == 12 {
			log("Pobrano ", unprocQuery, ", typu informacja klasyfikująca.")
			unprocQuery = RemoveSynonymes(unprocQuery)
			pocz := unprocQuery[:strings.Index(unprocQuery, "*be")]
			kon := unprocQuery[(strings.Index(unprocQuery, "*be") + 4):]
			//Trafne spostrzezenie dot. informacji
			fmt.Println(Type1Response(dbcontents, unprocQuery))
			//synonimy i zapamietanie
			SynonymeManagement(unprocQuery)
			addtodb(unprocQuery)
			AddClassification(kon, pocz)
		}
		if purpose == 2 {
			vout(dbcontents, "\n")
			log("Pobrano ", unprocQuery, ", typu pytanie.")
			placeofq := strings.Index(unprocQuery, "*q")
			qprefix := unprocQuery[:placeofq]
			qsuffix := unprocQuery[(placeofq + 2):]

			qprefix = strings.Trim(qprefix, " _\t\n!.")
			qsuffix = strings.Trim(qsuffix, " _\t\n!.")
			if qprefix == "()" {
				qprefix = ""
			}
			if qsuffix == "()" {
				qsuffix = ""
			}
			vout(string(placeofq), qprefix, "  ,  ", qsuffix)
			response := GrepIn(dbcontents, qprefix, qsuffix)
			vout("-----------------Info-koniec-----------------")
			fmt.Println(Type2Response(qprefix, qsuffix, response))
		}
		if purpose == 21 {
			log("Pobrano ", unprocQuery, ", typu pytanie tak/nie.")
			resp, isrly := GrepRly(unprocQuery)
			vout(resp, isrly)
			fmt.Println(Type21Response(resp, isrly))
		}
		if purpose == 22 {
			slicedQ := SliceAndTrim(unprocQuery)
			fmt.Println(Type22Response(GrepWhy(slicedQ[0], slicedQ[2])))
		}
		if purpose == 3 {
			slicedQ := SliceAndTrim(unprocQuery)
			slicedQ[1] = strings.Trim(slicedQ[1], " \n\t()/,.;!?")
			slicedQ[2] = strings.Trim(slicedQ[2], " \n\t()/,.;!?")
			vout("\""+slicedQ[1]+"\""+slicedQ[2]+"\"")
			aiRequests.Run(strings.Join([]string{slicedQ[1], slicedQ[2]}, " "), !(onRobot))
		}
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
	vout("qprefix:", qprefix, "\nqsuffix:", qsuffix, "\n")

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
	vout("GrepInPrefixMatch: ", itHasPrefix, "\n")
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
	vout("GrepInMatch:", answers, "\n")
	log("Dopasowałem odpowiedzi", answers)
	return answers
}
func GrepRly(query string) ([]string, bool) {
	trimmed, y2, y3 := SliceOnly(query)
	responses := make([]string, 32)
	respptr := 0
	var yesorno bool = false
	db := Scandb()
	for _, v := range db {
		if !(strings.Trim(v, " \n\t") == "") {
			trimv, y2v, y3v := SliceOnly(v)
			if (trimmed[0] == trimv[0])&&(trimmed[2] == trimv[2]) {
				if in(y2v, "*not")|| in(y3v, "*not") {
					if in(y2, "*not")|| in(y3, "*not") {
						responses[respptr] = v
						yesorno = true
						respptr += 1
					}
				} else {
					if !(in(y2, "*not"))||(!(in(y3, "*not"))) {
						responses[respptr] = v
						yesorno = true
						respptr += 1
					}
				}
			} else if (trimmed[2] == trimv[2]) && (slicesEq(y3, y3v)) {
				if in(y2v, "*not")|| in(y3v, "*not") {
					if in(y2, "*not")|| in(y3, "*not") {
						responses[respptr] = v
						yesorno = true
						respptr += 1
					}
				} else {
					if !(in(y2, "*not"))||(!(in(y3, "*not"))) {
						responses[respptr] = v
						yesorno = true
						respptr += 1
					}
				}
			}
		}
	}
	return responses, yesorno
}
func GetQuery() string {
	var inp string
	fmt.Printf(";> ")
	//źródło to konsola
	scnr := bufio.NewScanner(os.Stdin)
	//skanujemy i wynik do zmiennej
	scnr.Scan()
	inp = scnr.Text()
	if textformatting {
		inp = convert_new.Format(inp)
	}
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
		log("Pytanie", q, " tak/nie. Udzieliłe(a)ś odp. twierdzącej!")
		o = true
	} else {
		log("Pytanie", q, " tak/nie. Udzieliłe(a)ś odp. przeczącej!")
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
	log("Pobrałem dane z bazy danych...")
	return lines
}
func Querypurpose(query string) int8 {
	var querytype int8 = 1
	//jak zawiera *be to jest zdaniem klasyfikującym
	if strings.Contains(query, "*be") {
		querytype = 12
	}
	//jak zawiera *q to jest pytaniem
	if strings.Contains(query, "*q") {
		querytype = 2
	}
	//jeśli zawiera *why to jest pytaniem dlaczego
	if strings.Contains(query, "*why") {
		querytype = 22
	}
	//jak zawiera request to jest rozkazem
	if strings.Contains(query, "*r") {
		querytype = 3
	}
	//jeśli zawiera *rly to jest pytaniem tak/nie
	if strings.Contains(query, "*rly") {
		querytype = 21
	}
	//przywitania
	if strings.Contains(query, "*hi") {
		querytype = 101
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
	log("Dodałem", query, "do bazy danych")
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
	//log("Twoje zdanie zostało podzielone i skrócone.")
	return trimmed
}
func SliceOnly(query string) ([]string, []string, []string) {
	vout("slicing ", query)
	sliced := strings.Split(query, " ")
	trimmed := make([]string, 7)
	rawy2, rawy3 := "", ""
	if strings.Trim(query, " \n\t") != "" {
		var v string
		for i := 0; i < 3; i++ {
			v = sliced[i]
			trimmed[i] = v[:strings.Index(v, "(")]
		}
		rawy2 = sliced[1][(strings.Index(sliced[1], "(") + 1) : len(sliced[1])-1]
		rawy3 = sliced[2][(strings.Index(sliced[2], "(") + 1) : len(sliced[2])-1]

	}
	slicedy2 := strings.Split(rawy2, ",")
	slicedy3 := strings.Split(rawy3, ",")
	return trimmed, slicedy2, slicedy3
}
func Type1Response(db []string, query string) string {
	response := "Ok, didn't know that."
	sQuery := SliceAndTrim(RemoveSynonymes(query))
	for _, v := range db {
		v1, y2s, y3s := SliceOnly(v)
		y2, y3 := strings.Join(y2s, " "), strings.Join(y3s, " ")
		if v1[0] == sQuery[0] {
			response = "And I suppose it " + y2 + " " + v1[1] + " " + y3 + " " + v1[2] + " also..."
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
				v1, y2s, y3s := SliceOnly(v)
				y2, y3 := strings.Join(y2s, " "), strings.Join(y3s, " ")
				response = response + " " + y2 + " " + v1[1] + " " + y3 + " " + v1[2] + ",\n"
			}
		}
	} else if qp == "" {
		for _, v := range matching {
			if strings.Trim(v, " ") != "" {
				v1, y2s, _ := SliceOnly(v)
				y2 := strings.Join(y2s, " ")
				response = response + " " + v1[0] + " " + v1[1] + " " + y2 + ",\n"
			}
		}
		response = response + qs
		//response= strings.Join(matching, ", ")
	} else {
		response = qp
		for _, v := range matching {
			if strings.Trim(v, " ") != "" {
				v1, y2s, _ := SliceOnly(v)
				y2 := strings.Join(y2s, " ")
				response = response + " " + y2 + " " + v1[1] + ", "
			}
		}
		response = response + qs
		//response = matching[0]
	}
	return response
}
func Type21Response(causes []string, istrue bool) string {
	var response string
	if istrue {
		response = "Yes. "
		for _, v:=range causes {
			response += BracketsToEnglish(v)
		}
		
	} else {
		response = "No. Haven't heard 'bout that."
	}
	return response
}
func Type22Response (arguments []string) string {
	resp := "Because it's... "
	respl := ""
	vout(arguments)
	for _, v := range arguments {
		if strings.Trim(v, " \n\t") != "" {
			respl += "it's "+ v+", "
		}
	}
	return resp+respl
}
func DoubledInfoCheck() {
	db := Scandb()
	undoubleddb := make([]string, len(db))
	count := 0
	dblcnt := 0
	for _, v := range db {
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
			n2, err := f.Write(d2)
			errorcheck(err)
			defer vout("wrote %d bytes\n", n2)
		}
	}
	fmt.Printf("%v found...", (dblcnt - 1))
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
		vout(i)
		removed[i] = BaseWordOf(v[:strings.Index(v, "(")]) + v[strings.Index(v, "("):]
		vout(removed)
	}
	vout(sliced)
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
		vout(v, query)
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
			n2, err := f.Write(d2)
			errorcheck(err)
			defer vout("wrote %d bytes\n", n2)
		}
	}
}

func AddClassification(word, phrase string) {
	word = strings.Trim(word, " \t\n()")
	phrase = strings.Trim(phrase, " \t\n()")
	f, err := os.OpenFile("classif/"+word, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	errorcheck(err)
	//zakodowanie linii
	d2 := []byte(phrase + "\n")
	//piszemy
	n2, err := f.Write(d2)
	errorcheck(err)
	log("Sklasyfikowałem ", phrase, " jako ", word)
	defer fmt.Printf("wrote %d bytes\n", n2)
}
func GetClassification(word string) []string {
	classiflist, err := filepath.Glob("classif/*")
	errorcheck(err)
	//informatycznie odpowiednia długość wycinka
	lines := make([]string, 16383)
	if in(classiflist, "classif/"+word) {
		dat, err := ioutil.ReadFile("classif/" + word)
		errorcheck(err)
		data := string(dat)
		lines = strings.Split(data, "\n")
		return lines
	}
	return []string{""}
}
func DoubledClassificationCheck() {
	classiflist, err := filepath.Glob("classif/*")
	errorcheck(err)
	for _, vv := range classiflist {
		db := GetClassification(vv[(strings.Index(vv, "/") + 1):])
		undoubleddb := make([]string, len(db))
		count := 0
		dblcnt := 0
		for _, v := range db {
			if !(in(undoubleddb, strings.Trim(v, " \n\t"))) {
				undoubleddb[count] = strings.Trim(v, " \t\n")
				count += 1
			} else {
				dblcnt += 1
			}
		}
		f, err := os.Create(vv)
		defer f.Close()
		errorcheck(err)
		for _, v := range undoubleddb {
			if v != "" {
				d2 := []byte(v + "\n")
				n2, err := f.Write(d2)
				errorcheck(err)
				defer vout("wrote %d bytes\n", n2)
			}
		}
		fmt.Printf("%v, ", (dblcnt - 1))
	}
}
func ManageClassification(word, phrase string) {
	otherwords := GetClassification(phrase)
	AddClassification(word, phrase)
	for _, v := range otherwords {
		if v != "" {
			AddClassification(word, v)
		}
	}
}
func ListClassification() []string {
	classiflist, err := filepath.Glob("classif/*")
	errorcheck(err)
	retlist := make([]string, len(classiflist))
	for i, v := range classiflist {
		retlist[i] = v[strings.Index(v, "/")+1:]
	}
	return retlist
}

func GrepWhy(w1, w2 string) []string {
	causes := make([]string, 128)
	ccount := 0
	w2contents := GetClassification(w2)
	currentclass := make([]string, 64)
	for _, v:= range w2contents {
		currentclass = GetClassification(v)
		if in(currentclass, w1) {
			causes[ccount] = v
			ccount += 1
		}
	}
	return causes
}
func BracketsToEnglish(sentence string) string {
	rsentence := ""
	if strings.Trim(sentence, " \n\t()*!?.") != "" {
	xes, y2, y3 := SliceOnly(sentence)
	if xes[1]=="*be" {
		xes[1] = "is"
	}
	y2joined := strings.Join(y2, " ")
	y3joined := strings.Join(y3, " ")
	
	rsentence = xes[0]
	rsentence += " "+xes[1]+" "+y2joined
	rsentence += " the "+y3joined+" "+xes[2]+". "
}
	return rsentence
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
	t := time.Now()
	sz := fmt.Sprint(a...)
	s := fmt.Sprint("["+t.String()+"]", sz)
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	errorcheck(err)
	d2 := []byte(s + "\n")
	f.Write(d2)
	errorcheck(err)
}
func slicesEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
