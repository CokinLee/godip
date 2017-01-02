package godip

import (
	"fmt"
	"net/http"
)

var (
SALT_TEMPLATE = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN"
                      "http://www.w3.org/TR/html4/loose.dtd">
<html>
<head>
<title>
GoDIP Update Server
</title>
<meta name="salt" content="%s">
<meta name="time" content="%s">
<meta name="sign" content="%s">
</head>
<body>
<center>
<h2>
GoDIP Update Server
</h2>
Salt generated
</center>
</body>
</html>`
RETC_TEMPLATE = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN"
                      "http://www.w3.org/TR/html4/loose.dtd">
<html>
<head>
<title>
GoDIP Update Server
</title>
<meta name="retc" content="%s">
</head>
<body>
<center>
<h2>
GoDIP Update Server
</h2>
Successful update request
</center>
</body>
</html>`
)

func (c *Config) HttpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		q := r.URL.Query()
		//the "salt" from the first response ("salt=")
		salt := q.Get("salt")
		//the "time salt generated" value from the first response ("time=")
		//time_ := q.Get("time")
		//the "signature" from the first response ("sign=")
		sign := q.Get("sign")
		//the GnuDIP user name ("user=")
		user := q.Get("user")
		//the GnuDIP domain name ("domn=")
		domn := q.Get("domn")
		//the MD5 digested password created above ("pass=")
		pass := q.Get("pass")
		//the server "request code" ("reqc="):
		//    "0" - register the address passed with this request
		//    "1" - go offline
		//    "2" - register the address you see me at, and pass it back to me
		reqc := q.Get("reqc")
		//the IP address to be registered, if the request code is "0" ("addr=")
		addr := q.Get("addr")
		//REQUEST SALT
		if len(q) == 0 {
			signature = generateSign(r.Host)
			timeSalt, salt := generateSalt()
			fmt.Fprint(w, fmt.Sprintf(SALT_TEMPLATE, timeSalt, salt, signature))
		} else {
			switch reqc {
			case "0": //REQUEST UPDATE PROVIDING ADDRESS
				if c.auth(domn, salt, user, pass, sign, true) {
					c.Handler(&UpdateInformation{r, domn, user, r.RemoteAddr, RegisterUpdate})
					fmt.Fprint(w, fmt.Sprintf(RETC_TEMPLATE, "0"))
				} else {
					fmt.Fprint(w, fmt.Sprintf(RETC_TEMPLATE, "1"))
				}
			case "1": //OFFLINE REQUEST
				if c.auth(domn, salt, user, pass, sign, true) {
					c.Handler(&UpdateInformation{r, domn, user, r.RemoteAddr, GoOffline})
					fmt.Fprint(w, fmt.Sprintf(RETC_TEMPLATE, "2"))
				} else {
					fmt.Fprint(w, fmt.Sprintf(RETC_TEMPLATE, "1"))
				}
			case "2": //REQUEST UPDATE WITH ADDRESS SEEN BY SERVER
				if !c.auth(domn, salt, user, pass, sign, true) {
					c.Handler(&UpdateInformation{r, domn, user, r.RemoteAddr, RegisterBackIp})
					fmt.Fprint(w, fmt.Sprintf(RETC_TEMPLATE, `0">
<meta name="addr" content="`+addr))
				} else {
					fmt.Fprint(w, fmt.Sprintf(RETC_TEMPLATE, "1"))
				}
			}
		}
	}
	return
}
