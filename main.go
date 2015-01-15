// Package fetchall is a drop-in replacement for appengine/urlfetch
// to make requests exceeding the 32MB response size limit
package fetchall // import "timm.io/fetchall"

import (
	"fmt"
	"net"
	"net/http"

	"appengine"
	"appengine/socket"
)

func logDebugf(c appengine.Context, format string, args ...interface{}) {
	if appengine.IsDevAppServer() {
		c.Debugf(format, args...)
	}
}

func Client(c appengine.Context) *http.Client {
	var tr *http.Transport

	if appengine.IsDevAppServer() {
		// need to use `net` implementation with dev server, because of https://code.google.com/p/googleappengine/issues/detail?id=11076
		logDebugf(c, "using `net` implementation")
		tr = &http.Transport{
			Dial: func(network, hostPort string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(hostPort)
				if err != nil {
					return nil, err
				}

				addrs, err := net.LookupIP(host)
				if err != nil {
					return nil, err
				}
				logDebugf(c, "found addrs: %v", addrs)

				firstIP := addrs[0]
				var conn net.Conn
				if firstIP.To4() != nil {
					logDebugf(c, "first ip is ip4 %s", firstIP)
					conn, err = net.Dial(network, fmt.Sprintf("%s:%s", addrs[0], port))
				} else {
					// brackets for IPv6 addrs
					logDebugf(c, "first ip is ip6 %s", firstIP)
					conn, err = net.Dial(network, fmt.Sprintf("[%s]:%s", addrs[0], port))
				}

				if err != nil {
					return nil, err
				}
				logDebugf(c, "dialed, returning conn")

				return conn, nil
			},
		}
	} else {
		logDebugf(c, "using `appengine/socket` implementation")
		tr = &http.Transport{
			Dial: func(network, hostPort string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(hostPort)
				if err != nil {
					return nil, err
				}

				addrs, err := socket.LookupIP(c, host)
				if err != nil {
					return nil, err
				}
				logDebugf(c, "found addrs: %v", addrs)

				firstIP := addrs[0]
				var conn *socket.Conn
				if firstIP.To4() != nil {
					logDebugf(c, "first ip is ip4 %s", firstIP)
					conn, err = socket.Dial(c, network, fmt.Sprintf("%s:%s", addrs[0], port))
				} else {
					// brackets for IPv6 addrs
					logDebugf(c, "first ip is ip6 %s", firstIP)
					conn, err = socket.Dial(c, network, fmt.Sprintf("[%s]:%s", addrs[0], port))
				}

				if err != nil {
					return nil, err
				}
				logDebugf(c, "dialed, returning conn")

				return conn, nil
			},
		}
	}

	client := &http.Client{Transport: tr}

	return client
}
