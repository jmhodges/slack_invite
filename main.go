package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	subdomain      = flag.String("domain", "", "Slack subdomain to invite folks to. Required")
	email          = flag.String("email", "", "email address to send the Slack invite to. Required")
	weirdChannelId = flag.String("weirdChannelId", "", "weird channel id thing that can be used. Starts with a C")
	firstName      = flag.String("firstName", "", "first name of invitee")
	lastName       = flag.String("lastName", "", "last name of invitee")
)

func main() {
	token := strings.TrimSpace(os.Getenv("SLACK_TOKEN"))
	if token == "" {
		log.Fatalf("SLACK_TOKEN environment variable may not be blank")
	}
	if *subdomain == "" {
		flag.Usage() // TODO(jmhodges): improve error messaging
		os.Exit(2)
	}
	if !strings.HasSuffix(*subdomain, ".slack.com") {
		*subdomain = *subdomain + ".slack.com"
	}
	if *email == "" {
		flag.Usage() // TODO(jmhodges): improve error messaging
		os.Exit(2)
	}
	u := fmt.Sprintf(
		"https://%s/api/users.admin.invite?t=%d",
		*subdomain, time.Now().Unix(),
	)
	data := url.Values{}
	data.Set("email", *email)
	if *weirdChannelId != "" {
		data.Set("channels", *weirdChannelId)
	}
	data.Set("first_name", *firstName)
	data.Set("last_name", *lastName)
	data.Set("token", token)
	data.Set("set_active", "true")
	resp, err := http.PostForm(u, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read body: %s", err)
	}
	fmt.Println(string(b))
}
