package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Numbers struct {
	Numbers []int `json:"numbers"`
}

var window []int
var windowSize = 10

func handler(w http.ResponseWriter, r *http.Request) {
	numberType := r.URL.Path[len("/numbers/"):]
	url := "http://20.244.56.144/test/" + numberType

	client := http.Client{
		Timeout: 5 * time.Second, // Set a reasonable timeout
	}
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, "Error fetching data from external service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var numbers Numbers
	if err := json.Unmarshal(body, &numbers); err != nil {
		http.Error(w, "Error parsing response from external service", http.StatusInternalServerError)
		return
	}

	windowPrevState := make([]int, len(window))
	copy(windowPrevState, window)

	for _, number := range numbers.Numbers {
		if !contains(window, number) {
			if len(window) >= windowSize {
				window = window[1:]
			}
			window = append(window, number)
		}
	}

	avg := average(window)

	response := map[string]interface{}{
		"numbers":         numbers.Numbers,
		"windowPrevState": windowPrevState,
		"windowCurrState": window,
		"avg":             avg,
	}

	json.NewEncoder(w).Encode(response)
}

func contains(slice []int, item int) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func average(slice []int) float64 {
	if len(slice) == 0 {
		return 0.0 // Handle division by zero
	}
	sum := 0
	for _, v := range slice {
		sum += v
	}
	return float64(sum) / float64(len(slice))
}

func main() {
	http.HandleFunc("/numbers/", handler)
	http.ListenAndServe(":9876", nil)
}
