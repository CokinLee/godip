package godip

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func ListenAndServe(addr string, cfg *Config) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err.Error())
		}
		go func(c net.Conn) {
			defer c.Close()
			remoteHost := strings.Split(c.RemoteAddr().String(), ":")
			host := remoteHost[0]

			_, salt := generateSalt()
			c.Write([]byte(salt))

			rw := bufio.NewReader(c)
			buf := make([]byte, 1024)
			n, err := rw.Read(buf)
			if err != nil {
				log.Println(err.Error())
			}
			msg := strings.Split(string(buf[0:n]), ":")

			if len(msg) == 5 { // register
				//successful update
				if cfg.auth(msg[2], salt, msg[0], msg[1], "", false) {
					c.Write([]byte("0"))
					cfg.Handler(&UpdateInformation{nil, msg[2], msg[0], host, RegisterUpdate})
				} else { // invalid login
					c.Write([]byte("1"))
				}
			}
			if len(msg) == 4 {
				switch msg[3][:1] {
				case "0":
					// successful update
					if cfg.auth(msg[2], salt, msg[0], msg[1], "", false) {
						c.Write([]byte("0"))
						cfg.Handler(&UpdateInformation{nil, msg[2], msg[0], host, RegisterUpdate})
					} else { // invalid login
						c.Write([]byte("1"))
					}
				case "1":
					// successful offline
					if cfg.auth(msg[2], salt, msg[0], msg[1], "", false) {
						c.Write([]byte("2"))
						cfg.Handler(&UpdateInformation{nil, msg[2], msg[0], host, GoOffline})
					} else { // invalid login
						c.Write([]byte("1"))
					}
				case "2":
					// successful update and provides the address that was registered
					if cfg.auth(msg[2], salt, msg[0], msg[1], "", false) {
						c.Write([]byte(string("0:") + host))
						cfg.Handler(&UpdateInformation{nil, msg[2], msg[0], host, RegisterBackIp})
					} else { // invalid login
						c.Write([]byte("1"))
					}

				}
			}
		}(conn)
	}
}
