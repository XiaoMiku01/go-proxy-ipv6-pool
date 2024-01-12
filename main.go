package main

import (
	"crypto/rand"
	"log"
	"net"

	"flag"
	"net/http"
)

var cidr string

func main() {

	flag.StringVar(&cidr, "cidr", "", "ipv6 cidr")
	flag.Parse()

	// 监听socks5并服务
	go func() {
		err := socks5Server.ListenAndServe("tcp", "0.0.0.0:52122")
		if err != nil {
			log.Fatal(err)
		}

	}()
	// 在 8080 端口启动代理服务器
	log.Fatal(http.ListenAndServe(":52123", httpProxy))
}

func generateRandomIPv6(cidr string) (string, error) {
	// 解析CIDR
	_, ipv6Net, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}

	// 获取网络部分和掩码长度
	maskSize, _ := ipv6Net.Mask.Size()

	// 计算随机部分的长度
	randomPartLength := 128 - maskSize

	// 生成随机部分
	randomPart := make([]byte, randomPartLength/8)
	_, err = rand.Read(randomPart)
	if err != nil {
		return "", err
	}

	// 获取网络部分
	networkPart := ipv6Net.IP.To16()

	// 合并网络部分和随机部分
	for i := 0; i < len(randomPart); i++ {
		networkPart[16-len(randomPart)+i] = randomPart[i]
	}

	return networkPart.String(), nil
}
