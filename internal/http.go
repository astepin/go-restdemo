package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func HandleApiCalls() {
	http.HandleFunc("/random/mean", apiTest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func apiTest(wr http.ResponseWriter, r *http.Request) {

	requests, errR := strconv.ParseInt(r.URL.Query().Get("requests"), 10, 0)
	length, errL := strconv.ParseInt(r.URL.Query().Get("length"), 10, 0)

	if errR != nil || errL != nil || requests <= 0 || length <= 0 || requests > 100 || length > 50 {
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	url := "https://www.random.org/integers/?num=" + fmt.Sprint(length) + "&min=1&max=100&col=1&base=10&format=plain&rnd=new"

	var wg sync.WaitGroup
	var results []Result
	var lastError error

	wg.Add(int(requests))

	for idx := 1; idx <= int(requests); idx++ {
		go func() {
			err := getRandomNumbers(url, &wg, &results)
			if err != nil {
				lastError = err
				defer wg.Done()
			}
		}()
	}

	wg.Wait()

	if lastError != nil {
		sendError(wr, http.StatusInternalServerError, []byte(lastError.Error()))
		return
	}

	var stddevs []float64
	var allData []int

	for idx := 0; idx < len(results); idx++ {
		stddevs = append(stddevs, results[idx].Sd)
		allData = merge(allData, results[idx].Data)
	}
	allData = mergeSort(allData)

	results = append(results, Result{sd(stddevs), allData})

	fmt.Println(results)

	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)

	err := json.NewEncoder(wr).Encode(results)
	if err != nil {
		sendError(wr, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

}

func sendError(wr http.ResponseWriter, errorCode int, msg []byte) {

	wr.WriteHeader(errorCode)
	_, _ = wr.Write(msg)
}

func getRandomNumbers(url string, wg *sync.WaitGroup, results *[]Result) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.New("error getting random numbers")
	}

	// is response code 200
	if resp.StatusCode != http.StatusOK {
		return errors.New("random.org response code is not 200")
	}

	response, err := io.ReadAll(resp.Body)

	result := parseResult(string(response))

	if err != nil {
		return errors.New("error parsing response")
	}

	*results = append(*results, result)

	defer wg.Done()
	return nil
}
