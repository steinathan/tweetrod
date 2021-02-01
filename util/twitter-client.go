package util

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func StartTwitterClient() {

	config := oauth1.NewConfig("consumerKey", "consumerSecret")
	token := oauth1.NewToken("accessToken", "accessSecret")
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// twitter client
	client := twitter.NewClient(httpClient)

	//  Status Show
tweet, resp, err := client.Statuses.Show(585613041028431872, nil)

fmt.Println("tweet",tweet)
fmt.Println("resp",resp)
fmt.Println("err",err)
}
