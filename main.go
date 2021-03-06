package main //how basic Go programs begin

//imports all the functions going to be used in the code
import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	city := flag.String("city", "Pittsburgh", "Default city")
	//verbose gives enough detail to debug the program
	verbose := flag.Bool("v", false, "lots of info")
	//rawverbose gives all the details to debug the program
	rawVerbose := flag.Bool("rv", false, "lots and lots of info")
	live := flag.Bool("live", false, "whether or not to use live server")
	
	//applies the flags to find above to parsing command line parameters 
	flag.Parse()

	//sample appID by default
	appID := "b6907d289e10d714a6e88b30761fae22"
	//sample server by default
	server := "samples.openweathermap.org"
	//if live is turned on
	if *live {
		server = "api.openweathermap.org"
		appID = "246e1d08a3b875f4a75b7ca1b79fc7fe"
	}
	//%s is a placeholder value and the real values are displayed in order after the url
	url := fmt.Sprintf("http://%s/data/2.5/forecast?q=%s&mode=xml&appid=%s", server, *city, appID)
	fmt.Println(url)
	//creates an http Client which sends a request to the web server then the server sends a response back 
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
	//stores the data stream from the http server and stores it in b which is a buffer (array)
	io.Copy(&b, resp.Body)
	// spits out the rawest form of the data that can be found
	if *rawVerbose {
		fmt.Println(b.String())
	}
	w := Weather{}
	//parses the xml data and stores it in the weather struct 
	err = xml.Unmarshal(b.Bytes(), &w)
	if err != nil {
		panic(err)
	}
	fmt.Printf("location:%s, %s\n", w.Location.Name, w.Location.Country)
	if *verbose {
		//prints the raw data(%v) that was compiled 
		fmt.Printf("%v", w)
	}
	//searches for different precipitation and ranges of temperature 
	for _, t := range w.Forecast.Time {
		// searches for the word "snow" in weathermap.org
		if strings.Contains(t.Precipitation.Type, "snow") {
			fmt.Printf("%s - %s Snow on the way!\n", t.From, t.To)
		}
		// searches for the word "ice" in weathermap.org
		if strings.Contains(t.Precipitation.Type, "ice") {
			fmt.Printf("%s - %s Ice incoming\n", t.From, t.To)
		}
		//converts t.Temp.Value which is a string into a integer (float64)
		k, _ := strconv.ParseFloat(t.Temp.Value, 64)
		//formula for converting Kelvin from openweathermap.org to Fahrenheit
		f := ((9 / 5) * (k - 273)) + 32 
		//searches for temperature 32 degrees Fahrenheit and under
		if f <= 32 {
			fmt.Printf("%s - %s %d It'll be below freezing!\n", t.From, t.To, int(f))
		}
		//searches for temperatures above 32 degrees and below 40 degrees Fahrenheit
		if f > 32 && f <= 40 {
			fmt.Printf("%s - %s %d It'll be cold today!\n", t.From, t.To, int(f))
		}
		//searches for temperatures above 40 degrees and below 70 degrees Fahrenheit
		if f > 40 && f <= 70 {
			fmt.Printf("%s - %s %d It'll be warm today!\n", t.From, t.To, int(f))
		}
		//searches for temperatures greater than 70 degrees
		if f > 70 {
			fmt.Printf("%s- %s %d It'll be hot today!\n", t.From, t.To, int(f))
		}
	}
}
//http://samples.openweathermap.org/data/2.5/forecast?q=London,us&mode=xml&appid=b6907d289e10d714a6e88b30761fae22 
//tells unmarshal how to map the data structure above from the web server
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
	From          string        `xml:"from,attr"`
	To            string        `xml:"to,attr"`
	Symbol        Symbol        `xml:"symbol"`
	Precipitation Precipitation `xml:"precipitation"`
	Temp          Temp          `xml:"temperature"`
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

type Temp struct {
	Unit  string `xml:"unit,attr"`
	Value string `xml:"value,attr"`
	Min   string `xml:"min,attr"`
	Max   string `xml:"max,attr"`
}
