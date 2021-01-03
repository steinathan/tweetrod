package main

import (
	"flag"
	// "fmt"
	"math/rand"
	"runtime"
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

// BotPage ..
type BotPage struct {
	*rod.Page
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
		// devices.IPad,
		// devices.IPadPro,
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

	// wating for async results
	WaitForPageLoad(page)

	// now, get the users analytics tweet url && take some screenshot
	// then pass it to
	// var ch = make(chan string)
	// util.ProcessAnalytics(image, ch)
	// aUxl := CaptureScreenshot(page, "https://google.com")
	// fmt.Println(aUxl)
	time.Sleep(time.Second)
}

// CaptureScreenshot captures screenshot and returns file path, this file page will e fed directly to ProcessAnalytics for retieving the analytics data
func CaptureScreenshot(page *rod.Page, url string) string {
	var imgPath = "./imd.png"
	page.MustNavigate(url)

	WaitForPageLoad(page)

	page.MustScreenshot(imgPath)
	return imgPath
}

// WaitForPageLoad Waits until there's no more network connection (at the moment)
func WaitForPageLoad(page *rod.Page) {
	// an array to store long blocking events on twitter that stoped the next function from getting called
	// this `live_pipeline` used sockets long-polling and invalidates our code so we'll ignore it
	var excludes = make([]string, 0)

	util.Log("Blocking entity urls")
	excludes = []string{
		"https://api.twitter.com/live_pipeline/events",
		"https://twitter.com/i",
	}
	// wating for async results
	page.WaitRequestIdle(time.Duration(time.Second), []string{}, excludes)()

}
