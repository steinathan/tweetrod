package worker

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/otiai10/gosseract/v2"
)

// parses a string and returns the number equiv
func parseInt(s string) int {
	d, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		panic(err)
	}
	return int(d)
}

// CaptureScreenshot captures screenshot and returns file path, this file page will e fed directly to ProcessAnalytics for retieving the analytics data
func (w *RequestKind) CaptureScreenshot(page *rod.Page, url, username string) string {
	Log("Taking snapshot...")

	var imgPath = "./tmp/" + username + ".png"
	page.MustNavigate(url)

	WaitForPageLoad(page)

	Log("Waiting for page to be completely loaded...")

	page.MustScreenshot(imgPath)
	return imgPath
}

// ProcessAnalytics processes the analytics from a screenshot image
func (w *RequestKind) ProcessAnalytics(image string, ch chan AnalyticsKind) {
	Log("Loading OCR client...")

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

	Log("Getting impressions & engagements...")

	// compiled regular expressions for use
	var iRegex *regexp.Regexp = regexp.MustCompile("Impressions ([0-9]+)")
	var eRegex *regexp.Regexp = regexp.MustCompile("Total engagements ([0-9]+)")

	// We have our impressions and engagements
	impressKind := strings.Split(iRegex.FindString(text), " ")
	engageKind := strings.Split(eRegex.FindString(text), " ")

	if len(impressKind) == 1 || len(engageKind) < 2 {
		reason := fmt.Errorf("Could not get both impressions and engagements from provided shot: %v, %v", engageKind, impressKind)
		fmt.Println(reason)
		ch <- AnalyticsKind{
			IsValid:   false,
			TweetURL:  w.TweetURL,
			Username:  w.Username,
			CrawledAt: time.Now(),
		}
	} else {
		// process the structure
		var analytics = &AnalyticsKind{
			TweetURL:    w.TweetURL,
			Username:    w.Username,
			Engagements: parseInt(engageKind[2]),
			Impressions: parseInt(impressKind[1]),
			IsValid:     true,
			CrawledAt:   time.Now(),
		}

		ch <- *analytics

	}
	//:END

}
