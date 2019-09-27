package utils

import (
	"fmt"
	"net"
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

// GetIPv4Address iterates through the addresses expecting the format from
// func (ifi *net.Interface) Addrs() ([]net.Addr, error)
func GetIPv4Address(addresses []net.Addr) (string, error) {
	for _, addr := range addresses {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return "", err
		}
		ipv4 := ip.To4()
		if ipv4 == nil {
			continue
		}
		return ipv4.String(), nil
	}
	return "", fmt.Errorf("no addresses match")
}

// GetIPv6Address iterates through the addresses expecting the format from
// func (ifi *net.Interface) Addrs() ([]net.Addr, error) and returns the first
// non-link local address.
func GetIPv6Address(addresses []net.Addr) (string, error) {
	_, llNet, _ := net.ParseCIDR("fe80::/10")
	for _, addr := range addresses {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return "", err
		}
		if ip.To4() == nil && !llNet.Contains(ip) {
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("no addresses match")
}

// GetAddressForInterface looks for the network interface
// and returns the IPv4 address from the possible addresses.
func GetAddressForInterface(interfaceName string) (string, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		fmt.Printf("cannot find network interface %q: %v\n", interfaceName, err)
		return "", err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		fmt.Printf("cannot get addresses for network interface %q: %v\n", interfaceName, err)
		return "", err
	}
	return GetIPv4Address(addrs)
}

// GetV4OrV6AddressForInterface looks for the network interface
// and returns preferably the IPv4 address, and if it doesn't
// exists then IPv6 address.
func GetV4OrV6AddressForInterface(interfaceName string) (string, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		fmt.Printf("cannot find network interface %q: %v", interfaceName, err)
		return "", err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		fmt.Printf("cannot get addresses for network interface %q: %v", interfaceName, err)
		return "", err
	}
	if ip, err := GetIPv4Address(addrs); err == nil {
		return ip, nil
	}
	return GetIPv6Address(addrs)
}
