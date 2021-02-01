package util

import (
	// "encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

// AnalyticsKind the structure of our analytics
type AnalyticsKind struct {
	Engagements int `json:"engagements"`
	Impressions int `json:"impressions"`
}

// func main() {
// 	var ch = make(chan string) // AnalyticsKind as JSON.stringify(s)
// 	var img string = "./test.png"
// 	go ProcessAnalytics(img, ch)

// 	result := <-ch
// 	fmt.Println(result)

// }

// parses a string and returns the number equiv
func parseInt(s string) int {
	d, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		panic(err)
	}
	return int(d)
}

// ProcessAnalytics processes the analytics from a screenshot image
func ProcessAnalytics(image string, ch chan AnalyticsKind) {
	// ready our OCR
	client := gosseract.NewClient()
	defer client.Close()

	// set image & language
	client.SetImage(image)
	client.SetLanguage("eng")

	// gets the text
	text, err := client.Text()

	// handle error
	if err != nil {
		fmt.Println(err)
		return
	}

	// compiled regular expressions for use
	var iRegex *regexp.Regexp = regexp.MustCompile("Impressions ([0-9]+)")
	var eRegex *regexp.Regexp = regexp.MustCompile("Total engagements ([0-9]+)")

	// We have our impressions and engagements
	impressKind := strings.Split(iRegex.FindString(text), " ")
	engageKind := strings.Split(eRegex.FindString(text), " ")

	if len(impressKind) == 1 || len(engageKind) < 2 {
		reason := fmt.Errorf("Could not get both impressions and engagements from provided shot: %v, %v", engageKind, impressKind)
		panic(reason)
	}

	// process the structure
	var analytics = &AnalyticsKind{
		Engagements: parseInt(engageKind[2]),
		Impressions: parseInt(impressKind[1]),
	}

	// convert it to json
	// res, err := json.Marshal(analytics)
	// if err != nil {
	// 	panic(err)
	// }

	// send the converted json as a string to our channel
	ch <- *analytics

	//:END
}
