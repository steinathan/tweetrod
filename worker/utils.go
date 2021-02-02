package worker

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InputKind the structure of our analytics

// AnalyticsKind the structure of our analytics
type AnalyticsKind struct {
	IsValid     bool          `json:"isValid"`
	TweetURL    string        `json:"tweetURL"`
	Username    string        `json:"username"`
	Engagements int           `json:"engagements"`
	Impressions int           `json:"impressions"`
	CrawledAt   time.Duration `json:"crawledAt"`
}

// RequestKind ...
type RequestKind struct {
	Username string
	Password string
	TweetURL string
}

// HandleError ..
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

// WaitForPageLoad Waits until there's no more network connection (at the moment)
func WaitForPageLoad(page *rod.Page) {
	// an array to store long blocking events on twitter that stoped the next function from getting called
	// this `live_pipeline` used sockets long-polling and invalidates our code so we'll ignore it
	var excludes = make([]string, 0)

	Log("Blocking entity urls that may cause delay...")
	excludes = []string{
		"https://api.twitter.com/live_pipeline/events",
		"https://twitter.com/i",
	}
	// wating for async results
	page.WaitRequestIdle(time.Duration(time.Second), []string{}, excludes)()

}

// Log ..
func Log(format string, v ...interface{}) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.DisableCaller = true
	config.DisableStacktrace = true
	logger, _ := config.Build()
	logger.Info(fmt.Sprintf(format, v...))
}
