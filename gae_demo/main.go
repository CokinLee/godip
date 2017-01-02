package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"appengine"
	"appengine/urlfetch"

	"github.com/dgv/godip"
)

/* 
 * UPDATE APP ID AND ENVIRONMENT VARIABLES ON APP.YAML 
 */
 
func myHandler(u *godip.UpdateInformation) {
	if u.Reason == godip.RegisterUpdate {
		c := appengine.NewContext(u.HttpRequest)
		c.Infof("New IP update [%s@%s]: %s", u.User, u.Domain, u.IP)
		// Update DNS entry on Digital Ocean via API
		client := urlfetch.Client(c)
		req, err := http.NewRequest("PUT",
			fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records/%s", u.Domain, os.Getenv("DIGITALOCEAN_DOMAIN_ID")),
			strings.NewReader(fmt.Sprintf(`{"data": "%s"}`, u.IP)))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("DIGITALOCEAN_API_KEY")))
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			c.Infof("Error on digital ocean update: %s", err.Error())
			return
		}

		if resp.StatusCode == 200 {
			c.Infof("Record %s of %s successfully updated on digital ocean", os.Getenv("DIGITALOCEAN_DOMAIN_ID"), u.Domain)
		} else {
			c.Infof("Got %d from digital ocean's during record update of %s, check domain and ID", resp.StatusCode, u.Domain)
		}
	}
}

func init() {
	c := &godip.Config{
		Handler: myHandler,
		/*
		 * UPDATE WITH YOUR DOMAINS
		 */
		Domains: godip.Domains{"test1.com": {"router1": "test123"}},
	}

	http.HandleFunc("/gnudip/cgi-bin/gdipupdt.cgi", c.HttpHandler)
	http.HandleFunc("/cgi-bin/gdipupdt.cgi", c.HttpHandler)
}
