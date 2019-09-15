package utils

import (
	"fmt"
	"net"
	"os"
)

// GetLocalIP ...
func GetLocalIP() (ips []net.IP) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.To4())
			}
		}
	}
	return
}

// GetPrivaIP ...
func GetPrivaIP(ips []net.IP) (ip []string, err error) {

	for _, ipn := range ips {
		if IsPublicIP(ipn) {
			continue
		}
		ip = append(ip, ipn.To4().String())
	}
	if len(ip) == 0 {
		err = fmt.Errorf("NO PrivaIP")
	}
	return
}

// IsPublicIP ...
func IsPublicIP(IP net.IP) bool {
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

// FileExists ...
func FileExists(name string) (exists bool) {
	exists = true
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			exists = false
			return
		}
	}
	return
}
