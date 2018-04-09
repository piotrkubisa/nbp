package svc

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

const (
	// a - tabela kursów średnich walut obcych;
	// c - tabela kursów kupna i sprzedaży;
	avg                = "a"
	both               = "c"
	nbpAPI             = "http://www.nbp.pl/kursy/xml/"
	errCannotParseDate = "Couldn't parse the given date"
	errNbpAPIProblem   = "Couldn't get data from NBP API"
)

func getFormatedDate(date string) (string, error) {
	tmpDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", errors.New("Couldn't not parse the given date")
	}
	date = tmpDate.Format("060102")
	return date, nil
}

func getFormatedYear(date string) (string, error) {
	tmpDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", errors.New(errCannotParseDate)
	}

	fullYear := tmpDate.Format("2006")
	currentYear := time.Now().Format("2006")

	yr := ""
	if currentYear != fullYear {
		yr = fullYear
	}

	return yr, nil
}

// GetResourceLocation returns name of file that contains the currencies for given date
// We're searching the index which is just a txt file.
// E.g.: dir2015.txt - contains references to files that have currencies for 2015 year
// E.g.: dir.txt - contains references for the current year.
func GetResourceLocation(sDate string, sType string) (string, error) {

	yr, err := getFormatedYear(sDate)
	fmtDate, err := getFormatedDate(sDate)
	resp, err := http.Get(nbpAPI + "dir" + yr + ".txt")
	defer func() {
		err = resp.Body.Close()
	}()
	if err != nil {
		return "", errors.New(errNbpAPIProblem)
	}

	lns := bufio.NewReader(resp.Body)
	scn := bufio.NewScanner(lns)

	fLsts := []string{}

	for scn.Scan() {
		if strings.Contains(scn.Text(), fmtDate) {
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

// GetData fetches the currency/currencies in given file.
// There is one file for one date. E.g.: 2015-01-02 is a20123123.xml
func GetData(file string, code string) (Query, error) {

	resp, err := http.Get(nbpAPI + file + ".xml")
	defer func() {
		err = resp.Body.Close()
	}()
	if err != nil {
		return Query{}, errors.New("Couldn't not get data from NBP api")
	}

	var q Query
	currencyData, _ := ioutil.ReadAll(resp.Body)
	cData := bytes.NewReader(currencyData)

	decoder := xml.NewDecoder(cData)
	decoder.CharsetReader = charset.NewReaderLabel
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

// Query ...
type Query struct {
	FromData    string     `xml:"data_publikacji" json:"fromDate"`
	TableNumber string     `xml:"numer_tabeli" json:"tableNumber"`
	Currencies  []currency `xml:"pozycja" json:"currencies"`
}

type currency struct {
	Code    string `xml:"kod_waluty" json:"code"`
	Name    string `xml:"nazwa_waluty" json:"name"`
	Ratio   string `xml:"przelicznik" json:"ratio"`
	Average string `xml:"kurs_sredni" json:"average"`
	Buy     string `xml:"kurs_kupna" json:"buy"`
	Sell    string `xml:"kurs_sprzedazy" json:"sell"`
}

func Average(rDate, rCode string) (Query, error) {
	return makeCall(rDate, "avg", rCode)
}

func Both(rDate, rCode string) (Query, error) {
	return makeCall(rDate, "both", rCode)
}

func makeCall(rDate, rType, rCode string) (Query, error) {
	var res Query

	// Is the given date OK?
	date, err := time.Parse("2006-01-02", rDate)
	if err != nil {
		return res, fmt.Errorf("Given date is wrong. Use 'YYYY-MM-DD'")
	}

	// Disable future date
	sDate := date.Format("2006-01-02")
	if time.Now().Before(date) {
		return res, fmt.Errorf("Given date is wrong. Can't use future date")
	}

	// Disable date before 2002-01-02 -> first record in NBP
	minDate, _ := time.Parse("2006-01-02", "2002-01-02")
	if date.Before(minDate) {
		return res, fmt.Errorf("Given date is wrong. Min date is 2002-01-02")
	}

	// Get the file containing the the currency data
	f, err := GetResourceLocation(sDate, rType)
	if err != nil {
		return res, fmt.Errorf("There was some problem with your request")
	}

	// When the currency rate was not found for given date try to go back one day.
	// It's used to get currencies for holidays, or weekends
	prevData := date
	for f == "" {
		prevData = prevData.AddDate(0, 0, -1)
		f, err = GetResourceLocation(prevData.Format("2006-01-02"), rType)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if f == "" {
		return res, fmt.Errorf("Resource for given date was not found")
	}
	res, err = GetData(f, rCode)
	if err != nil {
		return res, fmt.Errorf("There was some problem with your request")
	}
	return res, nil
}
