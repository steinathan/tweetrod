package worker

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func process(page *rod.Page, inputs RequestKind) AnalyticsKind {
	imgPath := inputs.CaptureScreenshot(page, inputs.TweetURL, inputs.Username)

	var ch = make(chan AnalyticsKind) // AnalyticsKind
	var workerBot = RequestKind{
		TweetURL: inputs.TweetURL,
		Username: inputs.Username,
	}
	go workerBot.ProcessAnalytics(imgPath, ch)
	result := <-ch

	return result
}

// Bootstrap ..
func (inputs *RequestKind) Bootstrap() AnalyticsKind {

	var result AnalyticsKind
	var username, password string = inputs.Username, inputs.Password

	Log("Using %v goroutings because of CPU capabilities.", runtime.NumCPU())

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

	Log("Emulation server started with device: %v \n", device.Title)

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

		Log("Wrtiting cookies for the first time..")

		// start processing after login
		result = process(page, *inputs)

	}

	time.Sleep(time.Second)
	return result

}
