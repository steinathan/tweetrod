package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/navicstein/tweetrod/util"
)

// Analytics ..
type Analytics struct {
	impressions string
	engagements string
}

func main() {
	var username, password string
	flag.StringVar(&username, "username", "navicstein", "your twitter username")
	flag.StringVar(&password, "password", "842867", "your twitter password")
	var useProxy *bool = flag.Bool("use-proxy", false, "whether to use reverse proxy at http://127.0.0.1:8080")

	flag.Parse()

	util.Log("Using %v goroutings because of CPU capabilities.", runtime.NumCPU())
	util.Log("if this program keeps failing, re-start the command until a test is passed!")

	var newLauncher *launcher.Launcher = launcher.New()
	var instance *launcher.Launcher = newLauncher.
		Delete("use-mock-keychain") // delete flag "--use-mock-keychain"

	if *useProxy {
		instance = newLauncher.Proxy("127.0.0.1:8080")
		ch := make(chan error, 0)
		util.Log("Spawing proxy server at http://127.0.0.1:8080")
		go RunProxy(ch)
	}

	url := instance.MustLaunch()

	var browser = rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// constants that will be pushed to a config file and downloaded
	const (
		usernameSelector          = "#react-root > div > div > div.css-1dbjc4n.r-13qz1uu.r-417010 > main > div > div > div.css-1dbjc4n.r-13qz1uu > form > div > div:nth-child(6) > label > div > div.css-1dbjc4n.r-18u37iz.r-16y2uox.r-1wbh5a2.r-1udh08x.r-1inuy60.r-ou255f.r-vmopo1 > div > input"
		passwordSelector          = "#react-root > div > div > div.css-1dbjc4n.r-13qz1uu.r-417010 > main > div > div > div.css-1dbjc4n.r-13qz1uu > form > div > div:nth-child(7) > label > div > div.css-1dbjc4n.r-18u37iz.r-16y2uox.r-1wbh5a2.r-1udh08x.r-1inuy60.r-ou255f.r-vmopo1 > div > input"
		mustSeeSelectorAfterLogin = "#react-root > div > div > div.css-1dbjc4n.r-18u37iz.r-13qz1uu.r-417010 > main > div > div > div > div > div"
	)

	var twitterURL = "https://twitter.com/login?redirect_after_login=https://twitter.com/" + username

	// var baseURL = "https://twitter.com"
	// var redirectAfterLoginURL = baseURL + "/" + username

	// load the twitter login
	page := browser.MustPage(twitterURL)

	// Won't work on the Playground since the time is frozen.
	rand.Seed(time.Now().Unix())
	availDevices := []devices.Device{
		devices.IPhoneX,
		devices.IPad,
		devices.GalaxyNoteII,
		devices.IPadPro,
	}

	device := availDevices[rand.Int()%len(availDevices)]

	util.Log("Emulation server started with device: %v \n", device.Title)

	page.MustEmulate(device)

	// block ads
	proto.PageSetAdBlockingEnabled{
		Enabled: true,
	}.Call(page)

	// type in the username
	page.MustElement(usernameSelector).MustInput(username)

	// type in the password & hit ENTER
	page.MustElement(passwordSelector).MustInput(password).MustPress(input.Enter)

	// navigate to the profile
	// page.MustNavigate(redirectAfterLoginURL)

	// wating for async results
	WaitForPageLoad(page, true)

	// 	// attempt to grab some dummy text about the profile URL
	// 	text, err := page.MustElement(mustSeeSelectorAfterLogin).Text()
	// 	if err != nil {
	// 		fmt.Println("ERROR:", err)
	// 		return
	// 	} else if text == "" {
	// 		fmt.Println(`
	// ╔═╗┬─┐┬─┐┌─┐┬─┐┬
	// ║╣ ├┬┘├┬┘│ │├┬┘│
	// ╚═╝┴└─┴└─└─┘┴└─o`)
	// 		fmt.Println("ERROR: text test returned by crawler failed, try again", text)
	// 		return
	// 	}

	// 	// show the test to prove we're in a profile page
	// 	util.ShowTest(text)

	const analyticsEvaluationJS = `
	()=> {
		let links = [];
		const tweetLinks = Array.from(
		document.querySelectorAll("a")
		);
		tweetLinks.forEach((link) => {
		//if (link.href.includes("analytics"))
		 links.push(link.href);
		});
		return {
		...links,
		};
	}
	`

	page.MustElement(mustSeeSelectorAfterLogin).WaitLoad()
	// fmt.Printf("Wating for Asynchronous pages to return with DOM contents..\n\n")

	// WaitForPageLoad(page)

	// serde: evaluate the page and return found links
	// jsonPageContent := page.MustEval(analyticsEvaluationJS)
	// fmt.Println(jsonPageContent)
	tweetURL := "https://twitter.com/NavicsteinR/status/1340944589558083585/analytics"
	info := GetAnalytics(page, tweetURL)
	fmt.Println(tweetURL, info)

	// for _, v := range jsonPageContent.Map() {
	// 	tweetURL := v.String()
	// 	if tweetURL == "" {
	// 		panic("Sorry, there's no tweet activity here")
	// 	}
	// }

	fmt.Println(`
	// obviously, we'll need to harvest the analytics of a particular tweetID ..
	// and look it up 
	// let's assume we have a "self" tweetID that we want to see it's analytics..

	// since the user is logged in:
	// use the main twitter's API (instead of scraping), we'll loop through the list of tweets he has that corresponds with the ..
	// advertisers tweet we gave him, if a match is found, we'll quickly comb through the analytics and return the impressions and other ..
	// .. valid data in a CSV or and API or we'll hit the main servers "webhook" for each iteration in a loop if a successfull match is found

	// using this aproach solves the issues of pagination or infite scrolling while crawling for new tweet and makes sure that the overhead cost of threaded chromes ..
	// .. instances are reduced
	
	// In a nutshell here's my current implemetation from this point:
	// - harvest the users tweet until our match is found
	// - get the matched tweet and scrape the analytics
	// - return scraped data via goroutines as information to avoid deadlocks to other users
	`)

	fmt.Println("Sleeping for 1 minute then exiting...")
	time.Sleep(time.Hour)
}

// GetAnalytics returns the anaytics for a tweet link
func GetAnalytics(page *rod.Page, tweetURL string) Analytics {
	if !strings.Contains(tweetURL, "analytics") {
		util.Log("Can't harvest analytics for %v: no \"analytics\" in the url\n", tweetURL)
	}

	const analyticsEvaluationJS = `
	()=> {
		return {
			impressions: document.querySelector("body > div.ep-TweetPerformance.ep-Section > div.ep-ImpressionsSection > div > div > div.ep-MetricTopContainer > div.ep-MetricValue > span")?.innerText,
			engagements: document.querySelector("body > div.ep-TweetPerformance.ep-Section > div.ep-EngagementsSection > div.ep-Metric.ep-SubSection > div > div.ep-MetricTopContainer > div.ep-MetricValue > span")?.innerText
		};
	}
	`

	// goto the tweet URL & wait for network IDLE
	page.MustNavigate(tweetURL).MustWaitLoad()

	// wating for async results
	WaitForPageLoad(page, true)

	jsonAnalyticsContent := page.MustEval(analyticsEvaluationJS)
	fmt.Println("jsonAnalyticsContent:", jsonAnalyticsContent)

	for _, v := range jsonAnalyticsContent.Map() {
		util.Log("jsonAnalyticsContent %v", v.String())
		fmt.Println("jsonAnalyticsContentFMT", v.Raw())
	}

	if val, ok := jsonAnalyticsContent.Gets(); ok {
		fmt.Println("val:", val)
	}

	return Analytics{
		engagements: "todo!",
		impressions: "todo!",
	}
}

// WaitForPageLoad Waits until there's no more network connection (at the moment)
func WaitForPageLoad(page *rod.Page, exclude bool) {
	// an array to store long blocking events on twitter that stoped the next function from getting called
	// this `live_pipeline` used sockets long-polling and invalidates our code so we'll ignore it
	var excludes = make([]string, 0)

	if exclude {
		util.Log("Blocking entity urls")
		excludes = []string{
			"https://api.twitter.com/live_pipeline/events",
			"https://twitter.com/i",
		}
	}
	// wating for async results
	page.WaitRequestIdle(time.Duration(time.Second), []string{}, excludes)()

}
