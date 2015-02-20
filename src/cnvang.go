package main

import (
	"fmt"
	"strings"
)

//Funkcje sprawdzające zdanie
func GetLines(s string) []string {
	return strings.Split(s, " ")
}
func GetLine(s string, LineIndex int) string {
	return GetLines(s)[LineIndex]
}
func GetCommas(s string) []string {
	return strings.Split(s, ",")
}
func GetComma(s string, LineIndex int) string {
	return GetCommas(s)[LineIndex]
}

//Funkcje podrzędne
func dodnaw(slowo, co string) string {
	wynik := slowo + "(" + co + ")"
	return wynik
}

func usspoj(zdanie string, ile int) string {
	var slowo string = GetLine(zdanie, 0)
	var slowoy2 string = GetComma(zdanie, 0)
	if slowo == "A" {
		return strings.Replace(zdanie, "A ", "", ile)
	} else if slowo == "An" {
		return strings.Replace(zdanie, "An ", "", ile)
	} else if slowo == "The" {
		return strings.Replace(zdanie, "The ", "", ile)
	} else if slowoy2 == "a" {
		return strings.Replace(zdanie, "a", "", ile)
	} else if slowoy2 == "an" {
		return strings.Replace(zdanie, "an", "", ile)
	} else if slowoy2 == "the" {
		return strings.Replace(zdanie, "the", "", ile)
	}
	return zdanie
}

/*func usspoj(zdanie string, ile int) string {
	var ZdanieBezSpojnika, ZwrotGetLine, zgc string
	ZwrotGetLine = GetLine(zdanie, 0)
	zgc = GetComma(zdanie, 0)


	if ZwrotGetLine == "A" {
		ZdanieBezSpojnika = strings.Replace(zdanie, "A ", "", ile)
		return ZdanieBezSpojnika


	} else if zgc == "a" {
		zgc = strings.Replace(zdanie, "a", "", ile)
		return zgc


	} else if ZwrotGetLine == "An" {
		ZdanieBezSpojnika = strings.Replace(zdanie, "An ", "", ile)
		return ZdanieBezSpojnika


	} else if zgc == "an" {
		zgc = strings.Replace(zdanie, "an", "", ile)
		return zgc


	} else if ZwrotGetLine == "The" {
		ZdanieBezSpojnika = strings.Replace(zdanie, "The ", "", ile)
		return ZdanieBezSpojnika


	} else if zgc == "the" {
		zgc = strings.Replace(zdanie, "the", "", ile)
		return zgc
	}

	return ZdanieBezSpojnika
}*/

func ktoryspoj(zdanie string, dlugosc int) int {
	for dlugosc > 0 {
		if GetLine(zdanie, dlugosc) == "a" {
			return dlugosc
		} else if GetLine(zdanie, dlugosc) == "an" {
			return dlugosc
		} else if GetLine(zdanie, dlugosc) == "the" {
			return dlugosc
		}
		dlugosc--
	}
	return -1
}

func ustaly2(zdanie string, koniecy2 int) string {
	y2 := ""
	koniecy2_1 := koniecy2
	for koniecy2 > 1 {
		if koniecy2 != koniecy2_1 {
			y2 = GetLine(zdanie, koniecy2) + "," + y2
		} else if koniecy2 == koniecy2_1 {
			y2 = GetLine(zdanie, koniecy2) + y2
		}
		koniecy2--
	}
	return y2
}

func ustaly3(zdanie string, koniecy3, spojnik int) string {
	y3 := ""
	koniecy3_1 := koniecy3
	for koniecy3 > spojnik {
		if koniecy3 != koniecy3_1 {
			y3 = GetLine(zdanie, koniecy3) + "," + y3
		} else if koniecy3 == koniecy3_1 {
			y3 = GetLine(zdanie, koniecy3)
		}
		koniecy3--
	}
	return y3
}

//Funkcja GŁÓWNA
func format(zdanie string) string {
	zdanie = usspoj(zdanie, 1)
	var y2, x3, y3 string
	dlugosc := len(GetLines(zdanie)) - 1
	spojnik := ktoryspoj(zdanie, dlugosc)
	if spojnik == -1 {
		y2 = ustaly2(zdanie, dlugosc-2)
		x3 = GetLine(zdanie, dlugosc)
		y3 = GetLine(zdanie, dlugosc-1)
		if dlugosc == 2 {
			y3 = ""
		}
	} else if spojnik != -1 {
		y2 = ustaly2(zdanie, spojnik)
		y2 = usspoj(y2, -1)
		x3 = GetLine(zdanie, dlugosc)
		y3 = ustaly3(zdanie, dlugosc-1, spojnik)
		if dlugosc == 3 {
			y3 = ""
		}
	}
	wyn := dodnaw(GetLine(zdanie, 0), "") + " " + dodnaw(GetLine(zdanie, 1), y2) + " " + dodnaw(x3, y3)
	return wyn
}

//Pomocnicza funkcja wypisująca wynik funkcji głównej na ekranie
func main() {
	fmt.Println(format("You work on the project too long"))
	//fmt.Println(GetComma("the,a,boo", 2))
}
