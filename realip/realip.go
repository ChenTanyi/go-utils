package realip

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"syscall"
	"time"
)

var (
	// Ipv4Website the web or api to get real ipv4 address
	Ipv4Website = "http://ifconfig.me"

	// ErrUnknownAF the address family is neither AF_INET nor AF_INET6
	ErrUnknownAF = errors.New("Unknown address family")
	// ErrIPNotFound the ip with specific family was not found
	ErrIPNotFound = errors.New("IP not found")
)

var (
	client = &http.Client{}
)

// GetRealIP get real ip according to the address family
func GetRealIP(family int) ([]string, error) {
	return GetRealIPWithTimeout(family, 20*time.Second)
}

// GetRealIPWithTimeout get real ip according to the address family
func GetRealIPWithTimeout(family int, timeout time.Duration) ([]string, error) {
	if family == syscall.AF_INET {
		client.Timeout = timeout
		req, err := http.NewRequest("GET", Ipv4Website, nil)
		if err != nil {
			return nil, fmt.Errorf("Create request error: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Send request error: %v", err)
		}
		defer resp.Body.Close()
		addr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Read body error: %v", err)
		}
		return []string{string(addr)}, nil
	} else if family == syscall.AF_INET6 {
		ips := []string{}
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return nil, fmt.Errorf("Get interface addrs error: %v", err)
		}
		for _, addr := range addrs {
			if ip, ok := addr.(*net.IPNet); ok {
				if ip.IP.To4() == nil && ip.IP.IsGlobalUnicast() {
					ipv6 := ip.IP.String()
					ips = append(ips, ipv6)
				}
			}
		}
		if len(ips) == 0 {
			return nil, ErrIPNotFound
		}
		return ips, nil
	}
	return nil, ErrUnknownAF
}
