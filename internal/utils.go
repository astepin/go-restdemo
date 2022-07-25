package internal

import (
	"bufio"
	"log"
	"math"
	"strconv"
	"strings"
)

type Result struct {
	Sd   float64 `json:"stddev,omitempty"`
	Data []int   `json:"data,omitempty"`
}

func sd(num []float64) float64 {
	var sum float64
	var mean, sd float64

	if len(num) == 0 {
		return 0
	}

	if len(num) == 1 {
		return num[0]
	}

	for i := 0; i < len(num); i++ {
		sum += num[i]
	}
	mean = sum / float64(len(num))

	for j := 0; j < len(num); j++ {
		sd += math.Pow(num[j]-mean, 2)
	}
	sd = math.Sqrt(sd / float64(len(num)))

	return sd
}

func mergeSort(items []int) []int {
	if len(items) < 2 {
		return items
	}
	first := mergeSort(items[:len(items)/2])
	second := mergeSort(items[len(items)/2:])
	return merge(first, second)
}

func merge(a []int, b []int) []int {
	var final []int
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			final = append(final, a[i])
			i++
		} else {
			final = append(final, b[j])
			j++
		}
	}
	for ; i < len(a); i++ {
		final = append(final, a[i])
	}
	for ; j < len(b); j++ {
		final = append(final, b[j])
	}
	return final
}

func parseResult(s string) Result {

	var lines []int
	var sdLines []float64
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		intValue, err := strconv.Atoi(sc.Text())
		if err != nil {
			log.Fatal(err)
		}

		lines = append(lines, intValue)
		sdLines = append(sdLines, float64(intValue))

	}
	return Result{sd(sdLines), lines}
}
