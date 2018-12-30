package utility

import (
	"errors"
	"net"
	"os"
	"regexp"
	"sync"

	"github.com/fananchong/go-xserver/common"
)

var (
	ipinner     string
	ipouter     string
	onceipinner sync.Once
	onceipouter sync.Once
)

// GetIPInner : 获取内网 IP
func GetIPInner() string {
	onceipinner.Do(func() {
		switch common.XCONFIG.Network.IPType {
		case 0:
			ip, err := networkCard2IP(common.XCONFIG.Network.IPInner)
			if err != nil {
				common.XLOG.Errorln(err)
				os.Exit(1)
			}
			ipinner = ip
		default:
			ipinner = common.XCONFIG.Network.IPInner
		}
	})
	return ipinner
}

// GetIPOuter : 获取外网 IP
func GetIPOuter() string {
	onceipouter.Do(func() {
		switch common.XCONFIG.Network.IPType {
		case 0:
			ip, err := networkCard2IP(common.XCONFIG.Network.IPOuter)
			if err != nil {
				common.XLOG.Errorln(err)
				os.Exit(1)
			}
			ipouter = ip
		default:
			ipouter = common.XCONFIG.Network.IPOuter
		}
	})
	return ipouter
}

// GetIP : 根据类型获取IP
func GetIP(t common.IPType) string {
	switch t {
	case common.IPINNER:
		return GetIPInner()
	case common.IPOUTER:
		return GetIPOuter()
	default:
		panic("unknow ip type.")
	}
}

func networkCard2IP(name string) (string, error) {
	nic, err := net.InterfaceByName(name)
	if err != nil {
		return "", err
	}
	addresses, err := nic.Addrs()
	if err != nil {
		return "", err
	}
	r, _ := regexp.Compile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))`)
	for _, addr := range addresses {
		s := r.FindAllString(addr.String(), -1)
		if len(s) != 0 {
			return s[0], nil
		}
	}
	return "", errors.New("no find address. nic: " + name)
}

// GetIntranetListenPort : 获取服务器组内监听端口
func GetIntranetListenPort() int32 {
	return common.XCONFIG.Network.Port[1]
}
