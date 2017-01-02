package godip

import (
	"net/http"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"time"
)

type Domains map[string]map[string]string

const (
	//register the address passed with this request
	RegisterUpdate = iota
	// go offline
	GoOffline
	// register the address you see me at, and pass it back to me
	RegisterBackIp
)

type UpdateInformation struct {
	HttpRequest *http.Request
	Domain string
	User   string
	IP     string
	Reason int
}

type Config struct {
	Domains Domains
	Handler func(*UpdateInformation)
}

type Server struct {
	*Config
	Addr string
}

const SALT_SIGN = "8a49ecd1ca7b60dfbd3c0fad9ca16e49"

var signature string

func generateSign(host string) string {
	h := md5.New()
	io.WriteString(h, host+SALT_SIGN)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c *Config) auth(hostname, salt, user, password, sign string, isHttp bool) bool {
	if isHttp && sign != signature {
		log.Printf("invalid auth [%s@%s]: bad signature\n", user, hostname)
		return false
	}
	users, ok := c.Domains[hostname]
	if !ok {
		log.Printf("invalid auth [%s@%s]: domain not configured\n", user, hostname)
		return false
	}
	userPasswd, ok := users[user]
	if ok {
		h := md5.New()
		io.WriteString(h, userPasswd)
		userPasswd = fmt.Sprintf("%x", h.Sum(nil))
	} else {
		log.Printf("invalid auth [%s@%s]: user not found\n", user, hostname)
		return false
	}
	h := md5.New()
	io.WriteString(h, userPasswd+"."+salt)
	if fmt.Sprintf("%x", h.Sum(nil)) == password {
		return true
	}
	log.Printf("invalid auth [%s@%s]: password not match\n", user, hostname)
	return false
}

func generateSalt() (saltTime, salt string) {
	h := md5.New()
	saltTime = fmt.Sprintf("%v", time.Now().Unix())
	io.WriteString(h, saltTime)
	return saltTime, fmt.Sprintf("%x", h.Sum(nil))[:10]
}
