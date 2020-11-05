package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Response contains the full JSON from API
type Response struct {
}

type EnvVariables struct {
	Symbol string
	Ndays  int
}

// ApiKey is a Struct to describe the APIkey from a JSON file
type ApiKey struct {
	Apikey string
}

func getstocks(w http.ResponseWriter, r *http.Request) {

	// Get the API Key from ConfigMap
	APIKey, err := ioutil.ReadFile("config/config.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Grab and decode the ApiKey from JSON file
	apikey := string(APIKey)

	// Print to standard out for error checking
	fmt.Println(apikey)

	// Get Env variables
	Symbol := os.Getenv("Symbol")
	Ndays := os.Getenv("NDAYS")

	// Construct URL for API call
	url := ("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&apikey=" + apikey + "&symbol=" + Symbol)

	// Print to standard out for error checking
	fmt.Println(Symbol)
	fmt.Println(Ndays)
	fmt.Println(url)

	// GET request for API using APIKey and Symbol variable
	response, responseErr := http.Get(url)
	if responseErr != nil {
		fmt.Print(responseErr.Error())
		os.Exit(1)
	}

	defer response.Body.Close()

	// Decode response from JSON
	var data map[string]interface{}
	dataDecoder := json.NewDecoder(response.Body)
	decodeErr := dataDecoder.Decode(&data)

	// Error Checking against Null response
	if decodeErr != nil {
		log.Fatal(decodeErr)
	}

	// Grab RefreshedDate from Meta Date
	refreshedDate, refreshedErr := data["Meta Data"].(map[string]interface{})["3. Last Refreshed"]

	// Error Checking against Null value
	if !refreshedErr {
		log.Fatal(refreshedErr)
	}

	// Change format into string for use further on in code
	refreshedString := fmt.Sprintf("%v", refreshedDate)

	// Capture close price for Refreshed Date
	closePrice, ok := data["Time Series (Daily)"].(map[string]interface{})[refreshedString].(map[string]interface{})["4. close"]

	// Error check against lack of Inner Map
	if !ok {
		panic("inner map is not a map!")
	}

	// Change format into string for use further on in code
	closeString := fmt.Sprintf("%v", closePrice)

	fmt.Fprintln(w, Symbol, "data=[", closeString)

	// Get all the close values for each of the dates
	//for date, close := range data["Time Series (Daily)"].(map[string]interface{}) {
	//	fmt.Println("Date:", date, "Value:", close.(map[string]interface{})["4. close"])
	//}

}

func main() {
	http.HandleFunc("/", getstocks)
	http.ListenAndServe(":8080", nil)
}
