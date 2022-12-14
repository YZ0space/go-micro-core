package net

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	privateBlocks []*net.IPNet
)

func init() {
	// 本函数会返回IP地址和该IP所在的网络和掩码。
	// 例如，ParseCIDR("192.168.100.1/16")会返回IP地址192.168.100.1和IP网络192.168.0.0/16。
	for _, b := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} {
		if _, block, err := net.ParseCIDR(b); err == nil {
			privateBlocks = append(privateBlocks, block)
		}
	}
}

func ListenAddr(port string, fn func(string) (net.Listener, error)) (string, net.Listener, error) {
	ip, err := GetLocalIP()
	if err != nil {
		return "", nil, err
	}
	net, err := Listen(ip+port, fn)
	return ip, net, err
}

// Listen takes addr:portmin-portmax and binds to the first available port
// Example: Listen("localhost:5000-6000", fn)
func Listen(addr string, fn func(string) (net.Listener, error)) (net.Listener, error) {
	// host:port || host:min-max
	parts := strings.Split(addr, ":")

	//
	if len(parts) < 2 {
		return fn(addr)
	}

	// try to extract port range
	ports := strings.Split(parts[len(parts)-1], "-")

	// single port
	if len(ports) < 2 {
		return fn(addr)
	}

	// we have a port range

	// extract min port
	min, err := strconv.Atoi(ports[0])
	if err != nil {
		return nil, errors.New("unable to extract port range")
	}

	// extract max port
	max, err := strconv.Atoi(ports[1])
	if err != nil {
		return nil, errors.New("unable to extract port range")
	}

	// set host
	host := strings.Join(parts[:len(parts)-1], ":")

	// range the ports
	for port := min; port <= max; port++ {
		// try bind to host:port
		ln, err := fn(fmt.Sprintf("%s:%d", host, port))
		if err == nil {
			return ln, nil
		}

		// hit max port
		if port == max {
			return nil, err
		}
	}

	// why are we here?
	return nil, fmt.Errorf("unable to bind to %s", addr)
}

// GetLocalIP returns a real ip
func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("Failed to get interfaces! Err: %v", err)
	}
	var addrs []net.Addr
	for _, inter := range interfaces {

		if inter.Flags&net.FlagUp != 0 && inter.Flags&net.FlagBroadcast != 0 && !strings.Contains(inter.Name, "docker") { //过滤出开启且支持广播的网卡，排除docker虚拟网卡
			interAddrs, err := inter.Addrs()
			if err != nil {
				return "", fmt.Errorf("Failed to get interfaces addresses! Err: %v", err)
			}
			addrs = append(addrs, interAddrs...)
		}
	}

	var ipAddr []byte

	for _, rawAddr := range addrs {
		var ip net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP
		case *net.IPNet:
			ip = addr.IP
		default:
			continue
		}

		if ip.To4() == nil {
			continue
		}

		if !IsPrivateIP(ip.String()) {
			continue
		}

		ipAddr = ip
		break
	}

	if ipAddr == nil {
		return "", fmt.Errorf("No private IP address found, and explicit IP not provided")
	}

	return net.IP(ipAddr).String(), nil
}

func IsPrivateIP(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}
