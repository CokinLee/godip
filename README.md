GoDIP
=====

Go implementation of GnuDIP Dynamic DNS service. Using GnuDIP protocol you can easily configure the integration with your DNS server (BIND etc)
or external APIs (AWS-Route53, CloudFlare, DigitalOcean...), triggering changes by the IP update. Currently GnuDIP client are implemented in many network routers.

Check Google AppEngine demo folder to run a service for free.

### Library Usage 

On this example we run a NET/HTTP server at :3459/:8080 respectively, with the same configuration, updating a registry on Digital Ocean DNS's Server via API.

```go
package main

import (
	"log"
	"strings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"fmt"

	"github.com/dgv/godip"
)

func httpPut(url, apiKey, data string) (resp *http.Response, err error) {
		client := &http.Client{}
		req, err := http.NewRequest("PUT", url, strings.NewReader(data))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
		req.Header.Add("Content-Type", "application/json")
		resp, err = client.Do(req)
		defer resp.Body.Close()
		return
}

func myHandler(u *godip.UpdateInformation) {
	if u.Reason == godip.RegisterUpdate {
		log.Printf("New IP update [%s@%s]: %s", u.User, u.Domain, u.IP)
		// Update DNS entry on Digital Ocean via API
		resp, err := httpPut(fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records/%s", 
					 u.Domain, 
					 os.Getenv("DIGITALOCEAN_"+u.Domain+"_ID")),
				os.Getenv("DIGITALOCEAN_API_KEY"),
				fmt.Sprintf(`{"data": "%s"}`, u.IP))
		if err != nil {
			log.Printf("Error on digital ocean update: %s", err.Error())
			return
		}

		if resp.StatusCode == 200 {
			log.Printf("Record %s of %s successfully updated on digital ocean", 
						os.Getenv("DIGITALOCEAN_"+u.Domain+"_ID"), u.Domain)
		} else {
			log.Printf("Got %d from digital ocean's during record update of %s, 
						check domain and ID", resp.StatusCode, u.Domain)
		}
	}
}

func main() {
	// GoDIP common configuration
	c := &godip.Config{
		Handler: myHandler,
		// NOTE: ROUTER CONSIDERS FQND ON DOMAIN INPUT
		Domains: godip.Domains{"test.com": {"router1": "test123"},
			"test2.com": {"router1": "test123",
				"router2": "test321"},
			"test3.com": {"router3": "test123"}},
	}

	// GoDIP HTTP configuration
	http.HandleFunc("/gnudip/cgi-bin/gdipupdt.cgi", c.HttpHandler)
	http.HandleFunc("/cgi-bin/gdipupdt.cgi", c.HttpHandler)

	// start listeners
	go http.ListenAndServe(":8080", nil)
	go godip.ListenAndServe(":3459", c)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}

```
### To-do
- tests coverage
- godoc revision
- more dns api's integrations ex. CloudFlare, route53...

### Tested devices

- Huawei HG8245H

### License

Copyright 2017 Daniel G. Vargas

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
