package tools

import (
	"github.com/LCY2013/thinking-in-go/crontab/lib/errors"
	"net"
	"strings"
)

// GetLocalIP 获取本地网卡IP
func GetLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)

	// 获取所有网卡信息
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}

	// 获取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址：ipv4，ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); (isIpNet && !ipNet.IP.IsLoopback()) || strings.Contains(ipNet.IP.String(), "127.0.0.1") {
			// 跳过IPv6
			if ipNet.IP.To4() == nil {
				continue
			}
			ipv4 = ipNet.IP.String() // xxxx.xxxx.xxxx.xxxx
			return
		}
	}

	err = errors.ERR_NO_LOCAL_IP_FOUND

	return
}
