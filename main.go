package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

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

	"github.com/gofiber/fiber/v2"
)

// Analytics ..
type Analytics struct {
	impressions string
	engagements string
}

// IncomingRequestKind ...
type IncomingRequestKind struct {
	Username string
	Password string
	TweetURL string
}

// BotPage ..
type BotPage struct {
	*rod.Page
}

func startWorker(inputs *IncomingRequestKind) util.AnalyticsKind {

	var result util.AnalyticsKind
	var username, password string = inputs.Username, inputs.Password

	flag.Parse()

	util.Log("Using %v goroutings because of CPU capabilities.", runtime.NumCPU())

	var newLauncher *launcher.Launcher = launcher.New()
	var instance *launcher.Launcher = newLauncher.
		Delete("use-mock-keychain") // delete flag "--use-mock-keychain"

	url := instance.MustLaunch()

	var browser = rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// constants that will be pushed to a config file and downloaded
	const (
		usernameSelector          = "#react-root > div > div > div.css-1dbjc4n.r-13qz1uu.r-417010 > main > div > div > div.css-1dbjc4n.r-13qz1uu > form > div > div:nth-child(6) > label > div > div.css-1dbjc4n.r-18u37iz.r-16y2uox.r-1wbh5a2.r-l71dzp.r-1udh08x.r-1inuy60.r-ou255f.r-1b9bua6 > div > input"
		passwordSelector          = "#react-root > div > div > div.css-1dbjc4n.r-13qz1uu.r-417010 > main > div > div > div.css-1dbjc4n.r-13qz1uu > form > div > div:nth-child(7) > label > div > div.css-1dbjc4n.r-18u37iz.r-16y2uox.r-1wbh5a2.r-l71dzp.r-1udh08x.r-1inuy60.r-ou255f.r-1b9bua6 > div > input"
		mustSeeSelectorAfterLogin = "#react-root > div > div > div.css-1dbjc4n.r-18u37iz.r-13qz1uu.r-417010 > main > div > div > div > div > div"
	)

	var twitterURL = "https://twitter.com/login?redirect_after_login=" + inputs.TweetURL

	var userCookiePath string = "./cookies/" + username + ".json"

	var loadedCookiesStruct []*proto.NetworkCookieParam

	// read the users cookies
	userCookies, err := ioutil.ReadFile(userCookiePath)
	if err == nil {
		// only attempt to Unmarshal cookies when there's no error
		if err := json.Unmarshal(userCookies, &loadedCookiesStruct); err != nil {
			HandleError(err)
		}
	}

	// unmarshall the cookie into the struct
	var page *rod.Page = browser.MustPage(twitterURL)

	rand.Seed(time.Now().Unix())
	availDevices := []devices.Device{
		devices.IPhoneX,
	}

	device := availDevices[rand.Int()%len(availDevices)]

	util.Log("Emulation server started with device: %v \n", device.Title)

	page.MustEmulate(device)

	// load the users cookies for twitter
	if loadedCookiesStruct != nil {
		page.SetCookies(loadedCookiesStruct)

		// start processing after login loaded
		result = process(page, *inputs)

	} else {
		// else just log the user in normally

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

		cookies := page.MustCookies()
		cookiesByte, err := json.Marshal(cookies)
		HandleError(err)

		err = ioutil.WriteFile(userCookiePath, cookiesByte, 0644)
		HandleError(err)

		// start processing after login
		result = process(page, *inputs)

	}

	time.Sleep(time.Second)

	fmt.Println("RESULT:", result)
	return result

}

func process(page *rod.Page, inputs IncomingRequestKind) util.AnalyticsKind {
	imgPath := CaptureScreenshot(page, inputs.TweetURL, inputs.Username)

	var ch = make(chan util.AnalyticsKind) // AnalyticsKind
	go util.ProcessAnalytics(imgPath, ch)
	result := <-ch

	return result
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/:username", func(c *fiber.Ctx) error {
		var username, password, tweetURL = c.Params("username"), c.Query("password"), c.Query("tweetURL")
		// var result util.AnalyticsKind = startWorker(&IncomingRequestKind{
		// 	Username: username,
		// 	Password: password,
		// 	TweetURL: tweetURL,
		// })

		// data, _ := json.Marshal(result)
		// return c.SendString(string(data))
		msg := fmt.Sprintf("%s, %s, %s", username, password, tweetURL)
		return c.SendString(msg)
	})

	app.Listen(":3000")
}

// CaptureScreenshot captures screenshot and returns file path, this file page will e fed directly to ProcessAnalytics for retieving the analytics data
func CaptureScreenshot(page *rod.Page, url, username string) string {
	var imgPath = "./tmp/" + username + ".png"
	page.MustNavigate(url)

	WaitForPageLoad(page)
	fmt.Println("Waiting for page load..")

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

// HandleError ..
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
