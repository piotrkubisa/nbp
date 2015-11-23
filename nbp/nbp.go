package nbp

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
)

// a - tabela kursów średnich walut obcych;
// b - tabela kursów średnich walut niewymienialnych;
// c - tabela kursów kupna i sprzedaży;
// h - tabela kursów jednostek rozliczeniowych.
const (
	avg                = "a"
	both               = "c"
	nbpAPI             = "http://www.nbp.pl/kursy/xml/"
	errCannotParseDate = "Couldn't parse the given date"
	errNbpApiProblem   = "Couldn't get data from NBP API"
)

func GetResourceLocation(sDate string, sType string) (string, error) {

	pDate, err := time.Parse("2006-01-02", sDate)
	if err != nil {
		return "", errors.New(errCannotParseDate)
	}

	sDate = pDate.Format("060102")

	resp, err := http.Get(nbpAPI + "dir.txt")
	if err != nil {
		return "", errors.New(errNbpApiProblem)
	}

	defer resp.Body.Close()
	lns := bufio.NewReader(resp.Body)
	scn := bufio.NewScanner(lns)

	fLsts := []string{}

	for scn.Scan() {
		if strings.Contains(scn.Text(), sDate) {
			fLsts = append(fLsts, scn.Text())
		}
	}

	var resourceName string
	for _, entry := range fLsts {
		switch sType {
		case "avg":
			if string(entry[0]) == avg {
				resourceName = string(entry)
			}
		case "both":
			if string(entry[0]) == both {
				resourceName = string(entry)
			}
		default:
			resourceName = ""
		}

	}

	return resourceName, nil
}

func GetData(file string, code string) (Query, error) {

	resp, err := http.Get(nbpAPI + file + ".xml")
	if err != nil {
		return Query{}, errors.New("Couldn't not get data from NBP api")
	}
	defer resp.Body.Close()

	var q Query
	currencyData, _ := ioutil.ReadAll(resp.Body)
	cData := bytes.NewReader(currencyData)

	decoder := xml.NewDecoder(cData)
	decoder.CharsetReader = charset.NewReader
	err = decoder.Decode(&q)

	res := q
	res.Currencies = res.Currencies[:0]

	if code == "*" {
		return q, nil
	}

	codes := strings.Split(code, ",")
	for _, c := range q.Currencies {

		for _, cd := range codes {
			if cd == c.Code {
				res.Currencies = append(res.Currencies, c)
			}
		}

	}

	return res, nil

}

type Query struct {
	TableNumber string     `xml:"numer_tabeli" json:"tableNumber"`
	Currencies  []Currency `xml:"pozycja" json:"currencies"`
}

type Currency struct {
	Code    string `xml:"kod_waluty" json:"code"`
	Name    string `xml:"nazwa_waluty" json:"name"`
	Ratio   string `xml:"przelicznik" json:"ratio"`
	Average string `xml:"kurs_sredni" json:"average"`
	Buy     string `xml:"kurs_kupna" json:"buy"`
	Sell    string `xml:"kurs_sprzedazy" json:"sell"`
}
