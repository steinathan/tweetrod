package main

import (
	"encoding/json"

	"github.com/navicstein/tweetrod/worker"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Tweetrod... 👋!")
	})

	app.Post("/:username/:password", func(c *fiber.Ctx) error {

		var username, password, tweetURL = c.Params("username"), c.Params("password"), c.Query("tweetURL")

		worker.Log("Started logger for %s", username)

		var inputs = &worker.RequestKind{
			Username: username,
			Password: password,
			// append the analytics to the tweetTweetURL because the user ..
			// .. won't supply it to the main server,
			TweetURL: tweetURL + "/analytics",
		}

		var result = inputs.Bootstrap()

		data, _ := json.Marshal(result)

		worker.Log("Done evaluating %s", data)
		return c.SendString(string(data))
	})

	app.Listen(":3000")
}
