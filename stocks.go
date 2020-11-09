// This is my attempt at the test app
// Having looked at the raw data returned by the API, the logic I wanted to work through was as follows:
// 1. Get all the required variables from the system (env variables and from a configMap via k8s Secret)
// 2. Capture the full data output in a map.
// 3. Extract the Refreshed Date (not todays date because that wasn't consistantly returned).
// 4. Using the Refreshed date, work out the N valid weekdays backwards.
// 5. Create an slice of the cost data using the returned dates as keys.
// 6. Sum the total cost data.
// 7. Divide by the value of N to find the average.
// 8. Print the cost data in a line, in the required format with average at the end.
// 9. Write a httpHandler to listen on route '/' and return the output of the above data.
// Profit <- Not quite where we ended up.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Response contains the full JSON from API
type Response struct {
}

// Create struct for the Environment variables
type EnvVariables struct {
	Symbol string
	Ndays  int
}

// ApiKey is a Struct to describe the APIkey from a file
type ApiKey struct {
	Apikey string
}

const (
	layoutISO = "2006-01-02"
)

//func main() {
func getstocks(w http.ResponseWriter, r *http.Request) {

	// Get the API Key from file supplied via configMap/Secret
	APIKey, err := ioutil.ReadFile("config/config.txt")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Create APIKey variable
	apikey := string(APIKey)

	// Get Env variables
	Symbol := os.Getenv("Symbol")
	Ndays := os.Getenv("NDAYS")

	// Construct URL for API call
	url := ("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&apikey=" + apikey + "&symbol=" + Symbol)

	// Print to standard out for error checking
	fmt.Fprintln(w, apikey)
	fmt.Fprintln(w, Symbol)
	fmt.Fprintln(w, Ndays)
	fmt.Fprintln(w, url)

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
		fmt.Print(decodeErr.Error())
		os.Exit(1)
	}

	// Grab RefreshedDate from Meta Date
	refreshedDate, refreshedErr := data["Meta Data"].(map[string]interface{})["3. Last Refreshed"]

	// Change format from interface{} to string for use further on in code
	refreshedString := fmt.Sprintf("%v", refreshedDate)

	// Error Checking against Null value
	if refreshedErr != true {
		fmt.Print(refreshedErr)
		os.Exit(1)
	}

	// Convert date string into time.Date
	firstDate, _ := time.Parse(layoutISO, refreshedString)
	fmt.Fprintln(w, firstDate)

	// Convert to Integer
	NdaysInt, ndaysErr := strconv.Atoi(Ndays)

	if refreshedErr != true {
		fmt.Print(ndaysErr)
		os.Exit(1)
	}

	// Compute the date NDAYs before
	lastDate := firstDate.AddDate(0, 0, -NdaysInt)
	fmt.Fprintln(w, lastDate)

	// Create empty slice
	//var dates []string

	// For the range of firstDate to laastDate create slice of dates
	for rd := rangeDate(firstDate, lastDate); ; {
		date := rd()
		if date.IsZero() {
			break
		}
		fmt.Fprintln(w, date)
		//dates = append(dates, date.Format(layoutISO))
	}

	// Error checking, print dates
	//fmt.Print(dates)

	// Capture close price for Refreshed Date
	closePrice, ok := data["Time Series (Daily)"].(map[string]interface{})[refreshedString].(map[string]interface{})["4. close"]

	// Error check against lack of Inner Map
	if !ok {
		panic("Not an internal map")
	}

	// Change format into string for use further on in code
	closeString := fmt.Sprintf("%v", closePrice)

	fmt.Fprintln(w, Symbol, "data=[", closeString)

	// Get all the close values for each of the dates and return write to output
	for date, close := range data["Time Series (Daily)"].(map[string]interface{}) {
		fmt.Fprintln(w, date, close.(map[string]interface{})["4. close"])
	}

}

// Function returns a list of dates between start and end date range
func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

// Date returns the data as a date
func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// AddDate returns the date after NDAYS
func AddDate(years, months, days int) time.Time {
	return time.Date(years, time.Month(months), days, 0, 0, 0, 0, time.UTC)
}

func main() {
	http.HandleFunc("/", getstocks)
	http.ListenAndServe(":8080", nil)
}
