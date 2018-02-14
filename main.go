package main //how every Go program begins

//imports all the functions going to be used in the code
import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	city := flag.String("city", "Pittsburgh", "Default city")
	verbose := flag.Bool("v", false, "lots of info")
	rawVerbose := flag.Bool("rv", false, "lots and lots of info")
	live := flag.Bool("live", false, "whether or not to use live server")

	flag.Parse()

	//sample appID by default
	appID := "b6907d289e10d714a6e88b30761fae22"
	//sample server by default
	server := "samples.openweathermap.org"
	//if live is overridden
	if *live {
		server = "api.openweathermap.org"
		appID = "246e1d08a3b875f4a75b7ca1b79fc7fe"
	}
	url := fmt.Sprintf("http://%s/data/2.5/forecast?q=%s&mode=xml&appid=%s", server, *city, appID)
	fmt.Println(url)
	resp, err := http.Get(url)
	// couldnt talk to server
	if err != nil {
		panic(err)
	}
	// talked to server, didnt go well
	if resp.StatusCode >= 300 {
		panic(fmt.Errorf(resp.Status))
	}
	b := bytes.Buffer{}
	io.Copy(&b, resp.Body)
	// get loud
	if *rawVerbose {
		fmt.Println(b.String())
	}
	w := Weather{}
	err = xml.Unmarshal(b.Bytes(), &w)
	if err != nil {
		panic(err)
	}
	fmt.Printf("location:%s, %s\n", w.Location.Name, w.Location.Country)
	if *verbose {
		fmt.Printf("%v", w)
	}

	for _, t := range w.Forecast.Time {
		// searches for the word "snow" in weathermap.org
		if strings.Contains(t.Precipitation.Type, "snow") {
			fmt.Printf("%v - %v Snow on the way!\n", t.From, t.To)
		}
		// searches for the word "ice" in weathermap.org
		if strings.Contains(t.Precipitation.Type, "ice") {
			fmt.Printf("%v - %v Ice incoming\n", t.From, t.To)
		}
		if strings.Contains(t.Temperature.Value <= 273.15) {
			fmt.Println("%v - %v It'll be below freezing!\n", t.From, t.To)
		}
		if else strings.Contains(t.Temperature.Value >= 273.15 && <= 277.59) {
			fmt.Println("%v - %v It'll be cold today!\n", t.From, t.To)
		}
		if else strings.Contains(t.Temperature.Value >= 277.59 && <= 291.48) {
			fmt.Println("%v - %v It'll be warm today!\n", t.From, t.To)
		}
		if else strings.Contains(t.Temperature.Value >= 291.48) {
			fmt.Println("%v- %v It'll be hot today!\n", t.From, t.To)
		}
}

type Weather struct {
	XMLName  xml.Name `xml:"weatherdata"`
	Location Location `xml:"location"`
	Forecast Forecast `xml:"forecast"`
}

type Location struct {
	Name    string `xml:"name"`
	Country string `xml:"country"`
}

type Forecast struct {
	Time []Timeslot `xml:"time"`
}

type Timeslot struct {
	From   string `xml:"from,attr"`
	To     string `xml:"to,attr"`
	Symbol Symbol `xml:"symbol"`
	Precipitation Precipitation `xml:"precipitation"`
	Temp Temp `xml:"temperature"`
}

type Symbol struct {
	Number string `xml:"number,attr"`
	Name   string `xml:"name,attr"`
	Var    string `xml:"var,attr"`
}

type Precipitation struct {
	Unit  string `xml:"unit,attr"`
	Value string `xml:"value,attr"`
	Type  string `xml:"type,attr"`
}

type Temperature struct {
	Unit  string `xml:"unit, attr"`
	Value string `xml:"value, attr"`
	Min  string `xml:"min, attr"`
	Max  string `xml:"max, attr"`
}
